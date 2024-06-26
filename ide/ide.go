package ide

import (
	"fmt"
	"maps"
	"net/url"
	"sync"
	"testing"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/ide/mapper"
	"github.com/movio/gqlt/internal/config"
	"github.com/movio/gqlt/internal/eval"
	"github.com/movio/gqlt/internal/typecheck"
	"github.com/movio/gqlt/memosa"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/syn"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Logger interface {
	Debugf(fmt string, args ...any)
	Infof(fmt string, args ...any)
	Warnf(fmt string, args ...any)
	Errorf(fmt string, args ...any)
}

type IDE struct {
	ctx *memosa.Context
	// naive concurrency control
	// acquire read lock while snapshot is used,
	// acquire write lock to update
	lock sync.RWMutex
}

func New() *IDE {
	ctx := memosa.New()
	memosa.Set[sourcesInputQuery](ctx, Input{make(map[string]string)})
	return &IDE{ctx, sync.RWMutex{}}
}

func WithSnapshot[R any](ide *IDE, log Logger, f func(Snapshot) R) (R, error) {
	var err error
	s, cleanup := ide.snapshot(log)
	defer func() {
		defer cleanup()
		if r := recover(); r != nil {
			log.Errorf("panic: %v", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return f(s), err
}

func (ide *IDE) WithSnapshot(log Logger, f func(Snapshot)) error {
	_, err := WithSnapshot(ide, log, func(s Snapshot) struct{} {
		f(s)
		return struct{}{}
	})
	return err
}

func (ide *IDE) snapshot(log Logger) (Snapshot, func()) {
	ide.lock.RLock()
	return Snapshot{ide, log}, ide.lock.RUnlock
}

// A snapshot of the current state of the IDE.
// All ide operations are performed on a snapshot.
type Snapshot struct {
	ide *IDE
	log Logger
}

type Changes []Change

type Change interface {
	Apply(*Input)
}

type SetFileContent struct {
	Path    string
	Content string
}

var _ Change = SetFileContent{}

func (s SetFileContent) Apply(input *Input) {
	input.Sources[s.Path] = s.Content
}

func (ide *IDE) Apply(changes Changes) {
	ide.lock.Lock()
	defer ide.lock.Unlock()

	input := Input{ide.Sources()}
	for _, change := range changes {
		change.Apply(&input)
	}
	memosa.Set[sourcesInputQuery](ide.ctx, input)
}

func (ide *IDE) Sources() map[string]string {
	return maps.Clone(memosa.Get[sourcesInputQuery](ide.ctx).Sources)
}

func (ide *IDE) SetSchemas(schema *config.Schemas) {
	memosa.Set[schemasInputQuery](ide.ctx, schema)
}

func (ide *IDE) Schemas() *config.Schemas {
	return memosa.Get[schemasInputQuery](ide.ctx)
}

func (ide *IDE) Source(uri string) string {
	return ide.Sources()[uri]
}

func (s *Snapshot) Parse(uri string) Parsed[syn.File] {
	return memosa.Fetch[parseQuery](s.ide.ctx, parseKey{uri})
}

func (s *Snapshot) Mapper(uri string) *mapper.Mapper {
	return memosa.Fetch[mapperQuery](s.ide.ctx, mapperKey{uri})
}

func (s *Snapshot) Typecheck(uri string) typecheck.Info {
	return memosa.Fetch[typecheckQuery](s.ide.ctx, typecheckKey{uri})
}

type (
	schemaQuery struct{}
	schemaKey   struct{ URI string }
)

var _ memosa.Query[schemaKey, *syn.Schema] = schemaQuery{}

func (schemaQuery) Execute(ctx *memosa.Context, key schemaKey) *syn.Schema {
	path := lib.Must(url.Parse(key.URI)).Path
	schemas := memosa.Get[schemasInputQuery](ctx)
	return schemas.ForPath(path)
}

type (
	typecheckQuery struct{}
	typecheckKey   struct{ URI string }
)

var _ memosa.Query[typecheckKey, typecheck.Info] = typecheckQuery{}

func (typecheckQuery) Execute(ctx *memosa.Context, key typecheckKey) typecheck.Info {
	schemas := memosa.Fetch[schemaQuery](ctx, schemaKey(key))
	tcx := typecheck.New(schemas, &eval.Settings{})
	ast := memosa.Fetch[parseQuery](ctx, parseKey(key)).Ast
	return tcx.Check(ast)
}

type (
	mapperQuery struct{}
	mapperKey   struct{ Path string }
)

var _ memosa.Query[mapperKey, *mapper.Mapper] = mapperQuery{}

func (mapperQuery) Execute(ctx *memosa.Context, key mapperKey) *mapper.Mapper {
	return mapper.New(memosa.Get[sourcesInputQuery](ctx).Sources[key.Path])
}

func protoToPoint(mapper *mapper.Mapper, position protocol.Position) *ast.Point {
	line, col := int(position.Line), int(position.Character)
	point, err := mapper.ByteOffset(line, col)
	if err != nil {
		return nil
	}
	return &point
}

func posToProto(mapper *mapper.Mapper, position ast.HasPosition) protocol.Range {
	pos := position.Pos()
	startLine, startCol, err := mapper.LineAndColumn(pos.Start)
	if err != nil {
		panic(fmt.Sprintf("posToProto: %v: %v", pos, err))
	}

	endLine, endCol, err := mapper.LineAndColumn(pos.End)
	if err != nil {
		panic(fmt.Sprintf("posToProto: %v: %v", pos, err))
	}

	return protocol.Range{
		Start: protocol.Position{Line: uint32(startLine), Character: uint32(startCol)},
		End:   protocol.Position{Line: uint32(endLine), Character: uint32(endCol)},
	}
}

type logger struct {
	t testing.TB
}

func (l logger) Debugf(fmt string, args ...any) {
	l.t.Logf(fmt, args...)
}

func (l logger) Infof(fmt string, args ...any) {
	l.t.Logf(fmt, args...)
}

func (l logger) Warnf(fmt string, args ...any) {
	l.t.Logf(fmt, args...)
}

func (l logger) Errorf(fmt string, args ...any) {
	l.t.Errorf(fmt, args...)
}

func TestWith(t testing.TB, content string, f func(string, Snapshot)) {
	const path = "test.gqlt"

	ide := New()

	// working directory is `gqlt/ide`
	schemas, err := config.LoadSchemas("../")
	require.NoError(t, err)

	schema := schemas.ForPath(path)
	require.NotNil(t, schema, "should load gqlt/.graphqlrc.yaml")
	ide.SetSchemas(schemas)
	require.Equal(t, schema, ide.Schemas().ForPath(path))

	ide.Apply(Changes{SetFileContent{Path: path, Content: content}})
	require.NoError(t, ide.WithSnapshot(logger{t}, func(s Snapshot) {
		f(path, s)
	}))
}
