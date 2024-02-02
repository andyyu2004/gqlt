package internal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/andyyu2004/gqlt/gqlparser/formatter"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
	"github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/internal/slice"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/syn"
)

func (e *Executor) query(ctx context.Context, ecx *executionContext, expr *syn.OperationExpr) (any, error) {
	operation := expr.Operation
	for _, transform := range []transform{
		namespaceTransform{ecx.settings.namespace},
		variableTransform{schema: e.schema, scope: ecx.scope},
	} {
		operation = transform.transformOperation(operation)
	}

	usedFragments := map[string]struct{}{}
	observers := &validator.Events{}
	observers.OnFragmentSpread(func(_ *validator.Walker, fragmentSpread *syn.FragmentSpread) {
		usedFragments[fragmentSpread.Name.Value] = struct{}{}
	})

	validator.Walk(&syn.Schema{}, &syn.QueryDocument{Operations: []*syn.OperationDefinition{expr.Operation}}, observers)

	buf := formatOperation(operation)

	for fragmentName := range usedFragments {
		fragment, ok := ecx.scope.LookupFragment(fragmentName)
		// if fragment doesn't exist, we just ignore it. Will be caught by better validation later
		if ok {
			buf.WriteString("\n\n")
			buf.WriteString(fragment)
		}
	}

	// Pass our local variables directly also as graphql variables
	var data any
	req := Request{Query: buf.String(), Variables: ecx.scope.gqlVars()}
	errs, err := ecx.client.Request(ctx, req, &data)
	if err != nil {
		return nil, err
	}

	if len(errs) > 0 {
		return flatten(data), errs
	}

	return flatten(data), nil
}

func formatOperation(operation *syn.OperationDefinition) *bytes.Buffer {
	buf := bytes.NewBufferString("")

	f := formatter.NewFormatter(buf, formatter.WithIndent("  "))
	f.FormatQueryDocument(&syn.QueryDocument{
		Operations: []*syn.OperationDefinition{operation},
	})

	return buf
}

type transform interface {
	// Create a new operation definition that is a transformation of the given operation.
	// The original operation definition must not be mutated.
	// However, it is not required to make a deep copy.
	transformOperation(*syn.OperationDefinition) *syn.OperationDefinition
}

// add a namespace to the operation
// e.g. namespace ["foo", "bar"] will transform
// query { baz { qux } } to query { foo { bar { baz { qux } } } }
type namespaceTransform struct {
	namespace []string
}

func (t namespaceTransform) transformOperation(operation *syn.OperationDefinition) *syn.OperationDefinition {
	selectionSet := operation.SelectionSet
	// iterate in reverse order to build up the selection set from the inside out
	for i := len(t.namespace) - 1; i >= 0; i-- {
		selectionSet = syn.SelectionSet{
			&syn.Field{
				Alias:        lexer.Token{Value: t.namespace[i]},
				Name:         lexer.Token{Value: t.namespace[i]},
				SelectionSet: selectionSet,
			},
		}
	}

	return &syn.OperationDefinition{
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
	schema schema
	scope  *scope
	err    error
}

func (t variableTransform) transformOperation(operation *syn.OperationDefinition) *syn.OperationDefinition {
	// if there are variable defined then we pass the variables through as graphql variables
	if len(operation.VariableDefinitions) > 0 {
		return operation
	}

	// otherwise, we replace all variables with their values inline

	var topLevelType typename
	switch operation.Operation {
	case syn.Query:
		topLevelType = t.schema.QueryType
	case syn.Mutation:
		topLevelType = t.schema.MutationType
	case syn.Subscription:
		panic("subscriptions not supported")
	default:
		panic("unknown operation type")
	}

	return &syn.OperationDefinition{
		Operation:           operation.Operation,
		Name:                operation.Name,
		VariableDefinitions: nil, // drop all variable definitions as they have been "inlined"
		Directives:          operation.Directives,
		SelectionSet:        t.transformSelectionSet(topLevelType, operation.SelectionSet),
		Position:            operation.Position,
		Comment:             operation.Comment,
	}
}

func (t variableTransform) transformArgumentList(argTypes map[string]typename, argumentList syn.ArgumentList) syn.ArgumentList {
	return slice.Map(argumentList, func(argument *syn.Argument) *syn.Argument {
		argTy := argTypes[argument.Name.Value]
		return &syn.Argument{
			Name:     argument.Name,
			Value:    t.transformValue(argTy, argument.Value),
			Position: argument.Position,
			Comment:  argument.Comment,
		}
	})
}

func (t variableTransform) transformValue(expectedType typename, value *syn.Value) *syn.Value {
	ty := t.schema.Types[expectedType]
	switch value.Kind {
	case syn.Variable:
		//  unexpected children for variable value
		lib.Assert(len(value.Fields) == 0)
		val, ok := t.scope.Lookup(value.Raw)
		if !ok {
			t.err = fmt.Errorf("reference to undefined variable in graphql query: %s", value.Raw)
			return value
		}

		out := &syn.Value{Position: value.Position, Comment: value.Comment}
		switch val := val.(type) {
		case int:
			out.Kind = syn.IntValue
			out.Raw = fmt.Sprintf("%d", val)
		case float64:
			out.Kind = syn.FloatValue
			out.Raw = fmt.Sprintf("%f", val)
		case string:
			switch ty.Kind {
			case syn.Enum:
				out.Kind = syn.EnumValue
			default:
				out.Kind = syn.StringValue
			}
			out.Raw = val
		case bool:
			out.Kind = syn.BooleanValue
			out.Raw = fmt.Sprintf("%t", val)
		case nil:
			out.Kind = syn.NullValue
			out.Raw = "null"
		}

		return out
	default:
		return &syn.Value{
			Raw: value.Raw,
			Fields: slice.Map(value.Fields, func(child *syn.ChildValue) *syn.ChildValue {
				childTy := ty.InputFields[child.Name.Value]
				return &syn.ChildValue{
					Name:     child.Name,
					Value:    t.transformValue(childTy.Type, child.Value),
					Position: child.Position,
					Comment:  child.Comment,
				}
			}),
			Kind:               value.Kind,
			Position:           value.Position,
			Comment:            value.Comment,
			Definition:         value.Definition,
			VariableDefinition: value.VariableDefinition,
			ExpectedType:       value.ExpectedType,
		}
	}
}

func (t variableTransform) transformSelectionSet(ty typename, selectionSet syn.SelectionSet) syn.SelectionSet {
	return slice.Map(selectionSet, func(selection syn.Selection) syn.Selection {
		switch selection := selection.(type) {
		case *syn.Field:
			field := t.schema.Types[ty].Fields[selection.Name.Value]
			return &syn.Field{
				Alias:        selection.Alias,
				Name:         selection.Name,
				Arguments:    t.transformArgumentList(field.Args, selection.Arguments),
				Directives:   selection.Directives,
				SelectionSet: t.transformSelectionSet(field.Type, selection.SelectionSet),
				Position:     selection.Position,
				Comment:      selection.Comment,
			}
		case *syn.FragmentSpread:
			return &syn.FragmentSpread{
				Name:       selection.Name,
				Directives: selection.Directives,
				Position:   selection.Position,
				Comment:    selection.Comment,
			}
		case *syn.InlineFragment:
			return &syn.InlineFragment{
				TypeCondition: selection.TypeCondition,
				Directives:    selection.Directives,
				SelectionSet:  t.transformSelectionSet(typename(selection.TypeCondition.Value), selection.SelectionSet),
				Position:      selection.Position,
				Comment:       selection.Comment,
			}
		default:
			panic("unreachable")
		}
	})
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
