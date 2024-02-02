package ide

import (
	"maps"
	"sync"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/ide/mapper"
	"github.com/andyyu2004/gqlt/memosa"
	"github.com/andyyu2004/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type IDE struct {
	ctx *memosa.Context
	// naive concurrency control
	// acquire read lock while snapshot is used,
	// acquire write lock to update
	lock sync.RWMutex
}

func New() *IDE {
	ctx := memosa.New()
	memosa.Set[inputQuery](ctx, Input{make(map[string]string)})
	return &IDE{ctx, sync.RWMutex{}}
}

func (ide *IDE) Snapshot() (Snapshot, func()) {
	ide.lock.RLock()
	return Snapshot{ide}, ide.lock.RUnlock
}

// A snapshot of the current state of the IDE.
// All ide operations are performed on a snapshot.
type Snapshot struct {
	ide *IDE
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
	memosa.Set[inputQuery](ide.ctx, input)
}

func (ide *IDE) Sources() map[string]string {
	return maps.Clone(memosa.Fetch[inputQuery](ide.ctx, memosa.InputKey{}).Sources)
}

func (ide *IDE) Source(path string) string {
	return ide.Sources()[path]
}

func (s *Snapshot) Parse(path string) Parsed[syn.File] {
	return memosa.Fetch[parseQuery](s.ide.ctx, parseKey{path})
}

func (s *Snapshot) Mapper(path string) *mapper.Mapper {
	return memosa.Fetch[mapperQuery](s.ide.ctx, mapperKey{path})
}

type (
	mapperQuery struct{}
	mapperKey   struct{ Path string }
)

var _ memosa.Query[mapperKey, *mapper.Mapper] = mapperQuery{}

func (mapperQuery) Execute(ctx *memosa.Context, key mapperKey) *mapper.Mapper {
	return mapper.New(memosa.Fetch[inputQuery](ctx, memosa.InputKey{}).Sources[key.Path])
}

func posToProto(mapper *mapper.Mapper, position ast.HasPosition) protocol.Range {
	pos := position.Pos()
	startLine, startCol, err := mapper.LineAndColumn(pos.Start)
	if err != nil {
		panic(err)
	}

	endLine, endCol, err := mapper.LineAndColumn(pos.End)
	if err != nil {
		panic(err)
	}

	return protocol.Range{
		Start: protocol.Position{Line: uint32(startLine), Character: uint32(startCol)},
		End:   protocol.Position{Line: uint32(endLine), Character: uint32(endCol)},
	}
}
