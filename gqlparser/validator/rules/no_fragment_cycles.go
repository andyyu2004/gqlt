package validator

import (
	"fmt"
	"strings"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/andyyu2004/gqlt/gqlparser/validator"
	"github.com/andyyu2004/gqlt/syn"
)

func init() {
	AddRule("NoFragmentCycles", func(observers *Events, addError AddErrFunc) {
		visitedFrags := make(map[string]bool)

		observers.OnFragment(func(walker *Walker, fragment *syn.FragmentDefinition) {
			var spreadPath []*syn.FragmentSpread
			spreadPathIndexByName := make(map[string]int)

			var recursive func(fragment *syn.FragmentDefinition)
			recursive = func(fragment *syn.FragmentDefinition) {
				if visitedFrags[fragment.Name.Value] {
					return
				}

				visitedFrags[fragment.Name.Value] = true

				spreadNodes := getFragmentSpreads(fragment.SelectionSet)
				if len(spreadNodes) == 0 {
					return
				}
				spreadPathIndexByName[fragment.Name.Value] = len(spreadPath)

				for _, spreadNode := range spreadNodes {
					spreadName := spreadNode.Name.Value

					cycleIndex, ok := spreadPathIndexByName[spreadName]

					spreadPath = append(spreadPath, spreadNode)
					if !ok {
						spreadFragment := walker.Document.Fragments.ForName(spreadName)
						if spreadFragment != nil {
							recursive(spreadFragment)
						}
					} else {
						cyclePath := spreadPath[cycleIndex : len(spreadPath)-1]
						var fragmentNames []string
						for _, fs := range cyclePath {
							fragmentNames = append(fragmentNames, fmt.Sprintf(`"%s"`, fs.Name.Value))
						}
						var via string
						if len(fragmentNames) != 0 {
							via = fmt.Sprintf(" via %s", strings.Join(fragmentNames, ", "))
						}
						addError(
							Message(`Cannot spread fragment "%s" within itself%s.`, spreadName, via),
							At(spreadNode.Position),
						)
					}

					spreadPath = spreadPath[:len(spreadPath)-1]
				}

				delete(spreadPathIndexByName, fragment.Name.Value)
			}

			recursive(fragment)
		})
	})
}

func getFragmentSpreads(node syn.SelectionSet) []*syn.FragmentSpread {
	var spreads []*syn.FragmentSpread

	setsToVisit := []syn.SelectionSet{node}

	for len(setsToVisit) != 0 {
		set := setsToVisit[len(setsToVisit)-1]
		setsToVisit = setsToVisit[:len(setsToVisit)-1]

		for _, selection := range set {
			switch selection := selection.(type) {
			case *syn.FragmentSpread:
				spreads = append(spreads, selection)
			case *syn.Field:
				setsToVisit = append(setsToVisit, selection.SelectionSet)
			case *syn.InlineFragment:
				setsToVisit = append(setsToVisit, selection.SelectionSet)
			}
		}
	}

	return spreads
}
