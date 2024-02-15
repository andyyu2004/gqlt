package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("ScalarLeafs", func(observers *Events, addError AddErrFunc) {
		observers.OnField(func(walker *Walker, field *syn.Field) {
			if field.Definition == nil {
				return
			}

			fieldType := walker.Schema.Types[field.Definition.Type.Name()]
			if fieldType == nil {
				return
			}

			if fieldType.IsLeafType() && len(field.SelectionSet) > 0 {
				addError(
					Message(`Field "%s" must not have a selection since type "%s" has no subfields.`, field.Name.Value, fieldType.Name),
					At(field.Position),
				)
			}

			if !fieldType.IsLeafType() && len(field.SelectionSet) == 0 {
				addError(
					Message(`Field "%s" of type "%s" must have a selection of subfields.`, field.Name.Value, field.Definition.Type.String()),
					Suggestf(`"%s { ... }"`, field.Name.Value),
					At(field.Position),
				)
			}
		})
	})
}
