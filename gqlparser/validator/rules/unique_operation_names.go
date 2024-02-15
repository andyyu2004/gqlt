package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("UniqueOperationNames", func(observers *Events, addError AddErrFunc) {
		seen := map[string]bool{}

		observers.OnOperation(func(walker *Walker, operation *syn.OperationDefinition) {
			if seen[operation.Name.Value] {
				addError(
					Message(`There can be only one operation named "%s".`, operation.Name.Value),
					At(operation.Position),
				)
			}
			seen[operation.Name.Value] = true
		})
	})
}
