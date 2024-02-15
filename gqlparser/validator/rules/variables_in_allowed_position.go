package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("VariablesInAllowedPosition", func(observers *Events, addError AddErrFunc) {
		observers.OnValue(func(walker *Walker, value *syn.Value) {
			if value.Kind != syn.Variable || value.ExpectedType == nil || value.VariableDefinition == nil || walker.CurrentOperation == nil {
				return
			}

			tmp := *value.ExpectedType

			// todo: move me into walk
			// If there is a default non nullable types can be null
			if value.VariableDefinition.DefaultValue != nil && value.VariableDefinition.DefaultValue.Kind != syn.NullValue {
				if value.ExpectedType.NonNull {
					tmp.NonNull = false
				}
			}

			if !value.VariableDefinition.Type.IsCompatible(&tmp) {
				addError(
					Message(
						`Variable "%s" of type "%s" used in position expecting type "%s".`,
						value,
						value.VariableDefinition.Type.String(),
						value.ExpectedType.String(),
					),
					At(value.Position),
				)
			}
		})
	})
}
