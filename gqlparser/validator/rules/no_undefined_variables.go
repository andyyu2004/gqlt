package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("NoUndefinedVariables", func(observers *Events, addError AddErrFunc) {
		observers.OnValue(func(walker *Walker, value *syn.Value) {
			if walker.CurrentOperation == nil || value.Kind != syn.Variable || value.VariableDefinition != nil {
				return
			}

			if walker.CurrentOperation.Name.Value != "" {
				addError(
					Message(`Variable "%s" is not defined by operation "%s".`, value, walker.CurrentOperation.Name.Value),
					At(value.Position),
				)
			} else {
				addError(
					Message(`Variable "%s" is not defined.`, value),
					At(value.Position),
				)
			}
		})
	})
}
