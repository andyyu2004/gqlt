package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("UniqueVariableNames", func(observers *Events, addError AddErrFunc) {
		observers.OnOperation(func(walker *Walker, operation *syn.OperationDefinition) {
			seen := map[string]int{}
			for _, def := range operation.VariableDefinitions {
				// add the same error only once per a variable.
				if seen[def.Variable.Value] == 1 {
					addError(
						Message(`There can be only one variable named "$%s".`, def.Variable.Value),
						At(def.Position),
					)
				}
				seen[def.Variable.Value]++
			}
		})
	})
}
