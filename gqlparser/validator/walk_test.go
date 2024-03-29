package validator

import (
	"testing"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/parser"
	"github.com/movio/gqlt/syn"
	"github.com/stretchr/testify/require"
)

func TestWalker(t *testing.T) {
	schema, err := LoadSchema(Prelude, &ast.Source{Input: "type Query { name: String }\n schema { query: Query }"})
	require.Nil(t, err)
	query, err := parser.ParseQuery(&ast.Source{Input: "{ as: name }"})
	require.Nil(t, err)

	called := false
	observers := &Events{}
	observers.OnField(func(_ *Walker, field *syn.Field) {
		called = true

		require.Equal(t, "name", field.Name.Value)
		require.Equal(t, "as", field.Alias.Value)
		require.Equal(t, "name", field.Definition.Name)
		require.Equal(t, "Query", field.ObjectDefinition.Name)
	})

	Walk(schema, query, observers)

	require.True(t, called)
}

func TestWalkInlineFragment(t *testing.T) {
	schema, err := LoadSchema(Prelude, &ast.Source{Input: "type Query { name: String }\n schema { query: Query }"})
	require.Nil(t, err)
	query, err := parser.ParseQuery(&ast.Source{Input: "{ ... { name } }"})
	require.Nil(t, err)

	called := false
	observers := &Events{}
	observers.OnField(func(_ *Walker, field *syn.Field) {
		called = true

		require.Equal(t, "name", field.Name.Value)
		require.Equal(t, "name", field.Definition.Name)
		require.Equal(t, "Query", field.ObjectDefinition.Name)
	})

	Walk(schema, query, observers)

	require.True(t, called)
}
