package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("UniqueInputFieldNames", func(observers *Events, addError AddErrFunc) {
		observers.OnValue(func(walker *Walker, value *syn.Value) {
			if value.Kind != syn.ObjectValue {
				return
			}

			seen := map[string]bool{}
			for _, field := range value.Fields {
				if seen[field.Name.Value] {
					addError(
						Message(`There can be only one input field named "%s".`, field.Name.Value),
						At(field.Position),
					)
				}
				seen[field.Name.Value] = true
			}
		})
	})
}
