package validator

import (
	"fmt"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("KnownRootType", func(observers *Events, addError AddErrFunc) {
		// A query's root must be a valid type.  Surprisingly, this isn't
		// checked anywhere else!
		observers.OnOperation(func(walker *Walker, operation *syn.OperationDefinition) {
			var def *syn.Definition
			switch operation.Operation {
			case syn.Query, "":
				def = walker.Schema.Query
			case syn.Mutation:
				def = walker.Schema.Mutation
			case syn.Subscription:
				def = walker.Schema.Subscription
			default:
				// This shouldn't even parse; if it did we probably need to
				// update this switch block to add the new operation type.
				panic(fmt.Sprintf(`got unknown operation type "%s"`, operation.Operation))
			}
			if def == nil {
				addError(
					Message(`Schema does not support operation type "%s"`, operation.Operation),
					At(operation.Position))
			}
		})
	})
}
