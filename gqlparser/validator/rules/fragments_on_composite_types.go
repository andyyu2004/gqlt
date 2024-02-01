package validator

import (
	"fmt"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("FragmentsOnCompositeTypes", func(observers *Events, addError AddErrFunc) {
		observers.OnInlineFragment(func(walker *Walker, inlineFragment *syn.InlineFragment) {
			fragmentType := walker.Schema.Types[inlineFragment.TypeCondition.Value]
			if fragmentType == nil || fragmentType.IsCompositeType() {
				return
			}

			message := fmt.Sprintf(`Fragment cannot condition on non composite type "%s".`, inlineFragment.TypeCondition.Value)

			addError(
				Message(message),
				At(inlineFragment.Position),
			)
		})

		observers.OnFragment(func(walker *Walker, fragment *syn.FragmentDefinition) {
			if fragment.Definition == nil || fragment.TypeCondition.Value == "" || fragment.Definition.IsCompositeType() {
				return
			}

			message := fmt.Sprintf(`Fragment "%s" cannot condition on non composite type "%s".`, fragment.Name.Value, fragment.TypeCondition.Value)

			addError(
				Message(message),
				At(fragment.Position),
			)
		})
	})
}
