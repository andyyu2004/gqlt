package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/parser"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/andyyu2004/memosa"
)

type (
	parseQuery struct{}
	parseKey   struct{ Path string }
)

var _ memosa.Query[parseKey, syn.File] = parseQuery{}

func (parseQuery) Execute(ctx *memosa.Context, key parseKey) syn.File {
	files := memosa.Fetch[inputQuery](ctx, memosa.InputKey{})
	text := files.Sources[key.Path]
	src := &ast.Source{Name: key.Path, Input: text}

	parser, err := parser.New(src)
	if err != nil {
		// we fail hard if there's any lexing errors currently
		// TODO diagnostics
		return syn.File{}
	}

	ast, err := parser.Parse()
	if err != nil {
		// TODO diagnostics
	}

	return ast
}
