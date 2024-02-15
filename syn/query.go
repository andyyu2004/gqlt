package syn

import "github.com/movio/gqlt/gqlparser/lexer"

type Transform interface {
	// Create a new operation definition that is a transformation of the given operation.
	// The original operation definition must not be mutated.
	// However, it is not required to make a deep copy.
	TransformOperation(*OperationDefinition) *OperationDefinition
}

// add a namespace to the operation
// e.g. namespace ["foo", "bar"] will transform
// query { baz { qux } } to query { foo { bar { baz { qux } } } }
type NamespaceTransform struct {
	Namespace []string
}

func (t NamespaceTransform) TransformOperation(operation *OperationDefinition) *OperationDefinition {
	selectionSet := operation.SelectionSet
	// iterate in reverse order to build up the selection set from the inside out
	for i := len(t.Namespace) - 1; i >= 0; i-- {
		selectionSet = SelectionSet{
			&Field{
				Alias:        lexer.Token{Value: t.Namespace[i]},
				Name:         lexer.Token{Value: t.Namespace[i]},
				SelectionSet: selectionSet,
			},
		}
	}

	return &OperationDefinition{
		Operation:           operation.Operation,
		Name:                operation.Name,
		VariableDefinitions: operation.VariableDefinitions,
		Directives:          operation.Directives,
		SelectionSet:        selectionSet,
		Position:            operation.Position,
		Comment:             operation.Comment,
	}
}
