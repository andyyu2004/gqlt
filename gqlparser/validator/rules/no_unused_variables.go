package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("NoUnusedVariables", func(observers *Events, addError AddErrFunc) {
		observers.OnOperation(func(walker *Walker, operation *syn.OperationDefinition) {
			for _, varDef := range operation.VariableDefinitions {
				if varDef.Used {
					continue
				}

				if operation.Name != "" {
					addError(
						Message(`Variable "$%s" is never used in operation "%s".`, varDef.Variable, operation.Name),
						At(varDef.Position),
					)
				} else {
					addError(
						Message(`Variable "$%s" is never used.`, varDef.Variable),
						At(varDef.Position),
					)
				}
			}
		})
	})
}
