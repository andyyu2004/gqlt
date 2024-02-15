package validator

import (
	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("KnownDirectives", func(observers *Events, addError AddErrFunc) {
		type mayNotBeUsedDirective struct {
			Name   string
			Line   int
			Column int
		}
		seen := map[mayNotBeUsedDirective]bool{}
		observers.OnDirective(func(walker *Walker, directive *syn.Directive) {
			if directive.Definition == nil {
				addError(
					Message(`Unknown directive "@%s".`, directive.Name),
					At(directive.Position),
				)
				return
			}

			for _, loc := range directive.Definition.Locations {
				if loc == directive.Location {
					return
				}
			}

			// position must be exists if directive.Definition != nil
			tmp := mayNotBeUsedDirective{
				Name:   directive.Name,
				Line:   directive.Position.Line,
				Column: directive.Position.Column,
			}

			if !seen[tmp] {
				addError(
					Message(`Directive "@%s" may not be used on %s.`, directive.Name, directive.Location),
					At(directive.Position),
				)
				seen[tmp] = true
			}
		})
	})
}
