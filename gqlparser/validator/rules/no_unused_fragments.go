package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("NoUnusedFragments", func(observers *Events, addError AddErrFunc) {
		inFragmentDefinition := false
		fragmentNameUsed := make(map[string]bool)

		observers.OnFragmentSpread(func(walker *Walker, fragmentSpread *syn.FragmentSpread) {
			if !inFragmentDefinition {
				fragmentNameUsed[fragmentSpread.Name.Value] = true
			}
		})

		observers.OnFragment(func(walker *Walker, fragment *syn.FragmentDefinition) {
			inFragmentDefinition = true
			if !fragmentNameUsed[fragment.Name.Value] {
				addError(
					Message(`Fragment "%s" is never used.`, fragment.Name.Value),
					At(fragment.Position),
				)
			}
		})
	})
}
