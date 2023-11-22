package validator

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
)

func init() {
	AddRule("KnownFragmentNames", func(observers *Events, addError AddErrFunc) {
		observers.OnFragmentSpread(func(walker *Walker, fragmentSpread *ast.FragmentSpread) {
			if fragmentSpread.Definition == nil {
				addError(
					Message(`Unknown fragment "%s".`, fragmentSpread.Name),
					At(fragmentSpread.Position),
				)
			}
		})
	})
}
