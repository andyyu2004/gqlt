package validator

import (
	"fmt"
	"sort"
	"strings"

	//nolint:revive // Validator rules each use dot imports for convenience.
	. "github.com/movio/gqlt/gqlparser/validator"
	"github.com/movio/gqlt/syn"
)

func init() {
	AddRule("FieldsOnCorrectType", func(observers *Events, addError AddErrFunc) {
		observers.OnField(func(walker *Walker, field *syn.Field) {
			if field.ObjectDefinition == nil || field.Definition != nil {
				return
			}

			message := fmt.Sprintf(`Cannot query field "%s" on type "%s".`, field.Name.Value, field.ObjectDefinition.Name)

			if suggestedTypeNames := getSuggestedTypeNames(walker, field.ObjectDefinition, field.Name.Value); suggestedTypeNames != nil {
				message += " Did you mean to use an inline fragment on " + QuotedOrList(suggestedTypeNames...) + "?"
			} else if suggestedFieldNames := getSuggestedFieldNames(field.ObjectDefinition, field.Name.Value); suggestedFieldNames != nil {
				message += " Did you mean " + QuotedOrList(suggestedFieldNames...) + "?"
			}

			addError(
				Message(message),
				At(field.Position),
			)
		})
	})
}

// Go through all of the implementations of type, as well as the interfaces
// that they implement. If any of those types include the provided field,
// suggest them, sorted by how often the type is referenced,  starting
// with Interfaces.
func getSuggestedTypeNames(walker *Walker, parent *syn.Definition, name string) []string {
	if !parent.IsAbstractType() {
		return nil
	}

	possibleTypes := walker.Schema.GetPossibleTypes(parent)
	suggestedObjectTypes := make([]string, 0, len(possibleTypes))
	var suggestedInterfaceTypes []string
	interfaceUsageCount := map[string]int{}

	for _, possibleType := range possibleTypes {
		field := possibleType.Fields.ForName(name)
		if field == nil {
			continue
		}

		suggestedObjectTypes = append(suggestedObjectTypes, possibleType.Name)

		for _, possibleInterface := range possibleType.Interfaces {
			interfaceField := walker.Schema.Types[possibleInterface]
			if interfaceField != nil && interfaceField.Fields.ForName(name) != nil {
				if interfaceUsageCount[possibleInterface] == 0 {
					suggestedInterfaceTypes = append(suggestedInterfaceTypes, possibleInterface)
				}
				interfaceUsageCount[possibleInterface]++
			}
		}
	}

	suggestedTypes := append(suggestedInterfaceTypes, suggestedObjectTypes...)

	sort.SliceStable(suggestedTypes, func(i, j int) bool {
		typeA, typeB := suggestedTypes[i], suggestedTypes[j]
		diff := interfaceUsageCount[typeB] - interfaceUsageCount[typeA]
		if diff != 0 {
			return diff < 0
		}
		return strings.Compare(typeA, typeB) < 0
	})

	return suggestedTypes
}

// For the field name provided, determine if there are any similar field names
// that may be the result of a typo.
func getSuggestedFieldNames(parent *syn.Definition, name string) []string {
	if parent.Kind != syn.Object && parent.Kind != syn.Interface {
		return nil
	}

	possibleFieldNames := make([]string, 0, len(parent.Fields))
	for _, field := range parent.Fields {
		possibleFieldNames = append(possibleFieldNames, field.Name)
	}

	return SuggestionList(name, possibleFieldNames)
}
