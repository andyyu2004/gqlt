package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("KnownTypeNames", func(observers *Events, addError AddErrFunc) {
		observers.OnVariable(func(walker *Walker, variable *syn.VariableDefinition) {
			typeName := variable.Type.Name()
			typdef := walker.Schema.Types[typeName]
			if typdef != nil {
				return
			}

			addError(
				Message(`Unknown type "%s".`, typeName),
				At(variable.Position),
			)
		})

		observers.OnInlineFragment(func(walker *Walker, inlineFragment *syn.InlineFragment) {
			typedName := inlineFragment.TypeCondition.Value
			if typedName == "" {
				return
			}

			def := walker.Schema.Types[typedName]
			if def != nil {
				return
			}

			addError(
				Message(`Unknown type "%s".`, typedName),
				At(inlineFragment.Position),
			)
		})

		observers.OnFragment(func(walker *Walker, fragment *syn.FragmentDefinition) {
			typeName := fragment.TypeCondition.Value
			def := walker.Schema.Types[typeName]
			if def != nil {
				return
			}

			var possibleTypes []string
			for _, t := range walker.Schema.Types {
				possibleTypes = append(possibleTypes, t.Name)
			}

			addError(
				Message(`Unknown type "%s".`, typeName),
				SuggestListQuoted("Did you mean", typeName, possibleTypes),
				At(fragment.Position),
			)
		})
	})
}
