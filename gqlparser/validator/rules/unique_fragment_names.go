package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("UniqueFragmentNames", func(observers *Events, addError AddErrFunc) {
		seenFragments := map[string]bool{}

		observers.OnFragment(func(walker *Walker, fragment *syn.FragmentDefinition) {
			if seenFragments[fragment.Name.Value] {
				addError(
					Message(`There can be only one fragment named "%s".`, fragment.Name.Value),
					At(fragment.Position),
				)
			}
			seenFragments[fragment.Name.Value] = true
		})
	})
}
