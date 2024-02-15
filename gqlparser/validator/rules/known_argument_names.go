package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("KnownArgumentNames", func(observers *Events, addError AddErrFunc) {
		// A GraphQL field is only valid if all supplied arguments are defined by that field.
		observers.OnField(func(walker *Walker, field *syn.Field) {
			if field.Definition == nil || field.ObjectDefinition == nil {
				return
			}
			for _, arg := range field.Arguments {
				def := field.Definition.Arguments.ForName(arg.Name.Value)
				if def != nil {
					continue
				}

				var suggestions []string
				for _, argDef := range field.Definition.Arguments {
					suggestions = append(suggestions, argDef.Name)
				}

				addError(
					Message(`Unknown argument "%s" on field "%s.%s".`, arg.Name.Value, field.ObjectDefinition.Name, field.Name.Value),
					SuggestListQuoted("Did you mean", arg.Name.Value, suggestions),
					At(field.Position),
				)
			}
		})

		observers.OnDirective(func(walker *Walker, directive *syn.Directive) {
			if directive.Definition == nil {
				return
			}
			for _, arg := range directive.Arguments {
				def := directive.Definition.Arguments.ForName(arg.Name.Value)
				if def != nil {
					continue
				}

				var suggestions []string
				for _, argDef := range directive.Definition.Arguments {
					suggestions = append(suggestions, argDef.Name)
				}

				addError(
					Message(`Unknown argument "%s" on directive "@%s".`, arg.Name.Value, directive.Name),
					SuggestListQuoted("Did you mean", arg.Name.Value, suggestions),
					At(directive.Position),
				)
			}
		})
	})
}
