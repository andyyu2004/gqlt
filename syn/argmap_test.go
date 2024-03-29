package syn

import (
	"testing"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/stretchr/testify/require"
)

func TestArg2Map(t *testing.T) {
	defs := ArgumentDefinitionList{
		{Name: "A", Type: NamedType("String", ast.Position{}), DefaultValue: &Value{Kind: StringValue, Raw: "defaultA"}},
		{Name: "B", Type: NamedType("String", ast.Position{})},
	}

	t.Run("defaults", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{}, nil)
		require.Equal(t, "defaultA", args["A"])
		require.NotContains(t, args, "B")
	})

	t.Run("values", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{
			{Name: lexer.Token{Value: "A"}, Value: &Value{Kind: StringValue, Raw: "valA"}},
			{Name: lexer.Token{Value: "B"}, Value: &Value{Kind: StringValue, Raw: "valB"}},
		}, nil)
		require.Equal(t, "valA", args["A"])
		require.Equal(t, "valB", args["B"])
	})

	t.Run("nulls", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{
			{Name: lexer.Token{Value: "A"}, Value: &Value{Kind: NullValue}},
			{Name: lexer.Token{Value: "B"}, Value: &Value{Kind: NullValue}},
		}, nil)
		require.Equal(t, nil, args["A"])
		require.Equal(t, nil, args["B"])
		require.Contains(t, args, "A")
		require.Contains(t, args, "B")
	})

	t.Run("undefined variables", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{
			{Name: lexer.Token{Value: "A"}, Value: &Value{Kind: Variable, Raw: "VarA"}},
			{Name: lexer.Token{Value: "B"}, Value: &Value{Kind: Variable, Raw: "VarB"}},
		}, map[string]interface{}{})
		require.Equal(t, "defaultA", args["A"])
		require.NotContains(t, args, "B")
	})

	t.Run("nil variables", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{
			{Name: lexer.Token{Value: "A"}, Value: &Value{Kind: Variable, Raw: "VarA"}},
			{Name: lexer.Token{Value: "B"}, Value: &Value{Kind: Variable, Raw: "VarB"}},
		}, map[string]interface{}{
			"VarA": nil,
			"VarB": nil,
		})
		require.Equal(t, nil, args["A"])
		require.Equal(t, nil, args["B"])
		require.Contains(t, args, "A")
		require.Contains(t, args, "B")
	})

	t.Run("defined variables", func(t *testing.T) {
		args := arg2map(defs, ArgumentList{
			{Name: lexer.Token{Value: "A"}, Value: &Value{Kind: Variable, Raw: "VarA"}},
			{Name: lexer.Token{Value: "B"}, Value: &Value{Kind: Variable, Raw: "VarB"}},
		}, map[string]interface{}{
			"VarA": "varvalA",
			"VarB": "varvalB",
		})
		require.Equal(t, "varvalA", args["A"])
		require.Equal(t, "varvalB", args["B"])
	})
}
