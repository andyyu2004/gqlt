package validator

import (
	_ "embed"

	"github.com/movio/gqlt/gqlparser/ast"
)

//go:embed prelude.graphql
var preludeGraphql string

var Prelude = &ast.Source{
	Name:    "prelude.graphql",
	Input:   preludeGraphql,
	BuiltIn: true,
}
