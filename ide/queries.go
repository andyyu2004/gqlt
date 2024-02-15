package ide

import (
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/parser"
	"github.com/movio/gqlt/memosa"
	"github.com/movio/gqlt/syn"
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
	files := memosa.Fetch[sourcesInputQuery](ctx, memosa.InputKey{})
	text := files.Sources[key.URI]
	src := &ast.Source{Name: key.URI, Input: text}

	parser := parser.New(src)

	ast, err := parser.Parse()
	return Parsed[syn.File]{Ast: ast, Err: err}
}
