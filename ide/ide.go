package ide

import (
	"maps"
	"sync"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/andyyu2004/memosa"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type IDE struct {
	ctx  *memosa.Context
	once sync.Once
}

func New() *IDE {
	return &IDE{ctx: memosa.New()}
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
	ide.once.Do(func() { memosa.Set[inputQuery](ide.ctx, Input{make(map[string]string)}) })

	input := Input{maps.Clone(memosa.Fetch[inputQuery](ide.ctx, memosa.InputKey{}).Sources)}
	for _, change := range changes {
		change.Apply(&input)
	}
	memosa.Set[inputQuery](ide.ctx, input)
}

func (ide *IDE) Parse(path string) syn.File {
	return memosa.Fetch[parseQuery](ide.ctx, parseKey{path})
}

type Highlight struct {
	Pos       ast.Position
	TokenKind protocol.SemanticTokenType
}

type HighlightKind int

const (
	HighlightKindType HighlightKind = iota
)

func (ide *IDE) Highlight(path string) []Highlight {
	tree := ide.Parse(path)
	_ = tree
	return []Highlight{
		{ast.Position{Start: 0, End: 2, Line: 0, Column: 0}, protocol.SemanticTokenTypeType},
	}
}
