package gqlt

import (
	"bytes"
	"context"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/formatter"
	"github.com/andyyu2004/gqlt/syn"
)

func (e *Executor) query(ctx context.Context, ecx *executionContext, expr *syn.OperationExpr) (any, error) {
	operation := expr.Operation
	for _, transform := range []transform{
		namespaceTransform{ecx.settings.namespace},
		variableTransform{ecx.scope},
	} {
		operation = transform.transformOperation(operation)
	}

	query := formatOperation(operation)

	// Pass our local variables directly also as graphql variables
	var data any
	req := Request{Query: query, Variables: ecx.scope.gqlVars()}
	if err := e.client.Request(ctx, req, &data); err != nil {
		return nil, err
	}

	return flatten(data), nil
}

func formatOperation(operation *ast.OperationDefinition) string {
	buf := bytes.NewBufferString("")
	f := formatter.NewFormatter(buf, formatter.WithIndent("  "))
	f.FormatQueryDocument(&ast.QueryDocument{
		Operations: []*ast.OperationDefinition{operation},
	})
	return buf.String()
}

func mapSlice[T, U any](xs []T, f func(T) U) []U {
	ys := make([]U, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

type transform interface {
	// Create a new operation definition that is a transformation of the given operation.
	// The original operation definition must not be mutated.
	// However, it is not required to make a deep copy.
	transformOperation(*ast.OperationDefinition) *ast.OperationDefinition
}

// add a namespace to the operation
// e.g. namespace ["foo", "bar"] will transform
// query { baz { qux } } to query { foo { bar { baz { qux } } } }
type namespaceTransform struct {
	namespace []string
}

func (t namespaceTransform) transformOperation(operation *ast.OperationDefinition) *ast.OperationDefinition {
	selectionSet := operation.SelectionSet
	// iterate in reverse order to build up the selection set from the inside out
	for i := len(t.namespace) - 1; i >= 0; i-- {
		selectionSet = ast.SelectionSet{
			&ast.Field{Alias: t.namespace[i], Name: t.namespace[i], SelectionSet: selectionSet},
		}
	}

	return &ast.OperationDefinition{
		Operation:           operation.Operation,
		Name:                operation.Name,
		VariableDefinitions: operation.VariableDefinitions,
		Directives:          operation.Directives,
		SelectionSet:        selectionSet,
		Position:            operation.Position,
		Comment:             operation.Comment,
	}
}

// replace all variables with their values if no explicit parameter list
type variableTransform struct {
	scope *scope
}

func (t variableTransform) transformOperation(operation *ast.OperationDefinition) *ast.OperationDefinition {
	// if there are variable defined then we pass the variables through as graphql variables
	if len(operation.VariableDefinitions) > 0 {
		return operation
	}

	// otherwise, we replace all variables with their values inline
	// todo currently a noop

	return &ast.OperationDefinition{
		Operation:           operation.Operation,
		Name:                operation.Name,
		VariableDefinitions: nil,
		Directives:          operation.Directives,
		SelectionSet:        operation.SelectionSet,
		Position:            operation.Position,
		Comment:             operation.Comment,
	}
}

// flatten removes unnecessary nesting in a (hopefully) intuitive way from the graphql response
func flatten(data any) any {
	switch data := data.(type) {
	case map[string]any:
		if len(data) == 1 {
			for _, v := range data {
				return flatten(v)
			}
		}
		return data
	case []any:
		// recursively flatten elements of arrays
		xs := make([]any, len(data))
		for i, v := range data {
			xs[i] = flatten(v)
		}
		return xs
	default:
		return data
	}
}
