package syn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/parser"
	. "github.com/movio/gqlt/syn"
)

func TestQueryDocMethods(t *testing.T) {
	doc, err := parser.ParseQuery(&ast.Source{Input: `
		query Bob { foo { ...Frag } }
		fragment Frag on Foo {
			bar
		}
	`})

	require.Nil(t, err)
	t.Run("GetOperation", func(t *testing.T) {
		require.EqualValues(t, "Bob", doc.Operations.ForName("Bob").Name.Value)
		require.Nil(t, doc.Operations.ForName("Alice"))
	})

	t.Run("GetFragment", func(t *testing.T) {
		require.EqualValues(t, "Frag", doc.Fragments.ForName("Frag").Name.Value)
		require.Nil(t, doc.Fragments.ForName("Alice"))
	})
}

func TestNamedTypeCompatability(t *testing.T) {
	assert.True(t, NamedType("A", ast.Position{}).IsCompatible(NamedType("A", ast.Position{})))
	assert.False(t, NamedType("A", ast.Position{}).IsCompatible(NamedType("B", ast.Position{})))

	assert.True(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("A", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("B", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("B", ast.Position{}), ast.Position{})))

	assert.True(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("A", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("B", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("A", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("B", ast.Position{}), ast.Position{})))

	assert.True(t, NonNullNamedType("A", ast.Position{}).IsCompatible(NamedType("A", ast.Position{})))
	assert.False(t, NamedType("A", ast.Position{}).IsCompatible(NonNullNamedType("A", ast.Position{})))

	assert.True(t, NonNullListType(NamedType("String", ast.Position{}), ast.Position{}).IsCompatible(NonNullListType(NamedType("String", ast.Position{}), ast.Position{})))
	assert.True(t, NonNullListType(NamedType("String", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("String", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("String", ast.Position{}), ast.Position{}).IsCompatible(NonNullListType(NamedType("String", ast.Position{}), ast.Position{})))

	assert.True(t, ListType(NonNullNamedType("String", ast.Position{}), ast.Position{}).IsCompatible(ListType(NamedType("String", ast.Position{}), ast.Position{})))
	assert.False(t, ListType(NamedType("String", ast.Position{}), ast.Position{}).IsCompatible(ListType(NonNullNamedType("String", ast.Position{}), ast.Position{})))
}
