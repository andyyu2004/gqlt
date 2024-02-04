package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/parser"
	"github.com/andyyu2004/gqlt/memosa"
	"github.com/andyyu2004/gqlt/syn"
)

type (
	parseQuery struct{}
	parseKey   struct{ URI string }
)

var _ memosa.Query[parseKey, Parsed[syn.File]] = parseQuery{}

type Parsed[T any] struct {
	Ast syn.File
	Err error
}

func (parseQuery) Execute(ctx *memosa.Context, key parseKey) Parsed[syn.File] {
	files := memosa.Fetch[inputQuery](ctx, memosa.InputKey{})
	text := files.Sources[key.URI]
	src := &ast.Source{Name: key.URI, Input: text}

	parser := parser.New(src)

	ast, err := parser.Parse()
	return Parsed[syn.File]{Ast: ast, Err: err}
}
