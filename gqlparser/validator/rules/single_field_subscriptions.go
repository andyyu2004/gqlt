package validator

import (
	"strconv"
	"strings"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/syn"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
)

func init() {
	AddRule("SingleFieldSubscriptions", func(observers *Events, addError AddErrFunc) {
		observers.OnOperation(func(walker *Walker, operation *syn.OperationDefinition) {
			if walker.Schema.Subscription == nil || operation.Operation != syn.Subscription {
				return
			}

			fields := retrieveTopFieldNames(operation.SelectionSet)

			name := "Anonymous Subscription"
			if operation.Name.Value != "" {
				name = `Subscription ` + strconv.Quote(operation.Name.Value)
			}

			if len(fields) > 1 {
				addError(
					Message(`%s must select only one top level field.`, name),
					At(fields[1].position),
				)
			}

			for _, field := range fields {
				if strings.HasPrefix(field.name, "__") {
					addError(
						Message(`%s must not select an introspection top level field.`, name),
						At(field.position),
					)
				}
			}
		})
	})
}

type topField struct {
	name     string
	position ast.Position
}

func retrieveTopFieldNames(selectionSet syn.SelectionSet) []*topField {
	fields := []*topField{}
	inFragmentRecursive := map[string]bool{}
	var walk func(selectionSet syn.SelectionSet)
	walk = func(selectionSet syn.SelectionSet) {
		for _, selection := range selectionSet {
			switch selection := selection.(type) {
			case *syn.Field:
				fields = append(fields, &topField{
					name:     selection.Name.Value,
					position: selection.Pos(),
				})
			case *syn.InlineFragment:
				walk(selection.SelectionSet)
			case *syn.FragmentSpread:
				if selection.Definition == nil {
					return
				}
				fragment := selection.Definition.Name
				if !inFragmentRecursive[fragment.Value] {
					inFragmentRecursive[fragment.Value] = true
					walk(selection.Definition.SelectionSet)
				}
			}
		}
	}
	walk(selectionSet)

	seen := make(map[string]bool, len(fields))
	uniquedFields := make([]*topField, 0, len(fields))
	for _, field := range fields {
		if !seen[field.name] {
			uniquedFields = append(uniquedFields, field)
		}
		seen[field.name] = true
	}
	return uniquedFields
}
