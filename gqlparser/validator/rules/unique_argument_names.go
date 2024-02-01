package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("UniqueArgumentNames", func(observers *Events, addError AddErrFunc) {
		observers.OnField(func(walker *Walker, field *syn.Field) {
			checkUniqueArgs(field.Arguments, addError)
		})

		observers.OnDirective(func(walker *Walker, directive *syn.Directive) {
			checkUniqueArgs(directive.Arguments, addError)
		})
	})
}

func checkUniqueArgs(args syn.ArgumentList, addError AddErrFunc) {
	knownArgNames := map[string]int{}

	for _, arg := range args {
		if knownArgNames[arg.Name.Value] == 1 {
			addError(
				Message(`There can be only one argument named "%s".`, arg.Name.Value),
				At(arg.Position),
			)
		}

		knownArgNames[arg.Name.Value]++
	}
}
