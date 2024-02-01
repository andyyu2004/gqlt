package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("PossibleFragmentSpreads", func(observers *Events, addError AddErrFunc) {
		validate := func(walker *Walker, parentDef *syn.Definition, fragmentName string, emitError func()) {
			if parentDef == nil {
				return
			}

			var parentDefs []*syn.Definition
			switch parentDef.Kind {
			case syn.Object:
				parentDefs = []*syn.Definition{parentDef}
			case syn.Interface, syn.Union:
				parentDefs = walker.Schema.GetPossibleTypes(parentDef)
			default:
				return
			}

			fragmentDefType := walker.Schema.Types[fragmentName]
			if fragmentDefType == nil {
				return
			}
			if !fragmentDefType.IsCompositeType() {
				// checked by FragmentsOnCompositeTypes
				return
			}
			fragmentDefs := walker.Schema.GetPossibleTypes(fragmentDefType)

			for _, fragmentDef := range fragmentDefs {
				for _, parentDef := range parentDefs {
					if parentDef.Name == fragmentDef.Name {
						return
					}
				}
			}

			emitError()
		}

		observers.OnInlineFragment(func(walker *Walker, inlineFragment *syn.InlineFragment) {
			validate(walker, inlineFragment.ObjectDefinition, inlineFragment.TypeCondition.Value, func() {
				addError(
					Message(`Fragment cannot be spread here as objects of type "%s" can never be of type "%s".`, inlineFragment.ObjectDefinition.Name, inlineFragment.TypeCondition.Value),
					At(inlineFragment.Position),
				)
			})
		})

		observers.OnFragmentSpread(func(walker *Walker, fragmentSpread *syn.FragmentSpread) {
			if fragmentSpread.Definition == nil {
				return
			}
			validate(walker, fragmentSpread.ObjectDefinition, fragmentSpread.Definition.TypeCondition.Value, func() {
				addError(
					Message(`Fragment "%s" cannot be spread here as objects of type "%s" can never be of type "%s".`, fragmentSpread.Name.Value, fragmentSpread.ObjectDefinition.Name, fragmentSpread.Definition.TypeCondition.Value),
					At(fragmentSpread.Position),
				)
			})
		})
	})
}
