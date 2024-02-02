package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/memosa"
	"github.com/andyyu2004/gqlt/parser"
	"github.com/andyyu2004/gqlt/syn"
)

type (
	parseQuery struct{}
	parseKey   struct{ Path string }
)

var _ memosa.Query[parseKey, Parsed[syn.File]] = parseQuery{}

type Parsed[T any] struct {
	Ast syn.File
	Err error
}

func (parseQuery) Execute(ctx *memosa.Context, key parseKey) Parsed[syn.File] {
	files := memosa.Fetch[inputQuery](ctx, memosa.InputKey{})
	text := files.Sources[key.Path]
	src := &ast.Source{Name: key.Path, Input: text}

	parser, err := parser.New(src)
	if err != nil {
		// we fail hard if there's any lexing errors currently
		return Parsed[syn.File]{Err: err}
	}

	ast, err := parser.Parse()
	return Parsed[syn.File]{Ast: ast, Err: err}
}
