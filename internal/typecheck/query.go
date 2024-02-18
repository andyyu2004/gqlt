package typecheck

import (
	"fmt"

	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/internal/slice"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/syn"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Placeholder type to represent recursive types.
// We cannot eagerly expand out the graphql type because it may be recursive (and therefore infinite) without this type.
// This type should not leak out of this file, it should be expanded as much as necessary as dictated by the query.
type namedType struct {
	Name lex.Token
}

var _ Ty = namedType{}

func (ty namedType) isTy() {}

func (ty namedType) String() string {
	return ty.Name.Value
}

func (tcx *typechecker) queryExpr(expr *syn.QueryExpr) Ty {
	if tcx.schema == nil {
		// Impossible to typecheck a query without a schema.
		// It would be possible to approximate it except for the fact we have no way of telling when a field is a list or not.
		// Could do it based off whether the field name is plural or not, but that's only a heuristic and disgusting.
		return Any{}
	}

	ty := tcx.inferQueryType(syn.NamespaceTransform{
		Namespace: tcx.settings.Namespace(),
	}.TransformOperation(expr.Operation))
	return flatten(ty)
}

func (tcx *typechecker) inferQueryType(operation *syn.OperationDefinition) Ty {
	lib.Assert(tcx.schema != nil, "schema must not be nil to use this method")

	var def *syn.Definition
	switch operation.Operation {
	case syn.Query:
		def = tcx.schema.Query
	case syn.Mutation:
		def = tcx.schema.Mutation
	case syn.Subscription:
		def = tcx.schema.Subscription
	default:
		panic(fmt.Sprintf("unknown operation type %v", operation.Operation))
	}

	if def == nil {
		tcx.error(operation, fmt.Sprintf("schema does not define a %v type", operation.Operation))
	}

	// top level type is always an object
	ty := tcx.defToTy(def).(Object)
	return tcx.inferSelectionSetType(ty, operation.SelectionSet)
}

func (tcx *typechecker) inferSelectionSetType(ty Object, selectionSet syn.SelectionSet) Ty {
	outTy := Object{Fields: orderedmap.New[string, Ty]()}
	for _, selection := range selectionSet {
		switch selection := selection.(type) {
		case *syn.Field:
			fieldTy, ok := ty.Fields.Get(selection.Name.Value)
			if !ok {
				fieldTy = tcx.error(selection, fmt.Sprintf("field '%v' does not exist on type '%v'", selection.Name.Value, ty))
			}

			fieldTy = tcx.expand(fieldTy)
			if isErr(fieldTy) {
				outTy.Fields.Set(selection.Alias.Value, fieldTy)
				continue
			}

			switch f := fieldTy.(type) {
			case Any:
			case Object:
				if len(selection.SelectionSet) == 0 {
					fieldTy = tcx.error(selection, fmt.Sprintf("field '%s' of type '%v' must have a selection of subfields", selection.Name.Value, f))
				} else {
					fieldTy = tcx.inferSelectionSetType(f, selection.SelectionSet)
				}
			case List:
				switch elem := tcx.expand(f.Elem).(type) {
				case Object:
					if len(selection.SelectionSet) == 0 {
						fieldTy = List{Elem: tcx.error(selection, fmt.Sprintf("field '%s' of type '%v' must have a selection of subfields", selection.Name.Value, f))}
					} else {
						fieldTy = List{Elem: tcx.inferSelectionSetType(elem, selection.SelectionSet)}
					}
				case namedType:
					panic("should be expanded")
				}
			case namedType:
				panic("should be expanded")
			default:
				if len(selection.SelectionSet) > 0 {
					fieldTy = tcx.error(selection, fmt.Sprintf("cannot query field '%s' on type '%v'", selection.Name.Value, fieldTy))
				}
			}

			outTy.Fields.Set(selection.Alias.Value, fieldTy)
		case *syn.FragmentSpread:
			fragment, ok := tcx.fragments[selection.Name.Value]
			if !ok {
				outTy.Fields.Set(selection.Name.Value, tcx.error(selection, fmt.Sprintf("fragment '%v' not defined", selection.Name.Value)))
				continue
			}
			return tcx.inferSelectionSetType(ty, fragment.SelectionSet)
		case *syn.InlineFragment:
			return Any{}
		}
	}

	return outTy
}

func (tcx *typechecker) graphqlTyToTy(ty *syn.Type) Ty {
	return tcx.expand(tcx.graphqlTyToTy1(ty))
}

func (tcx *typechecker) graphqlTyToTy1(ty *syn.Type) Ty {
	lib.Assert(ty != nil, "type must not be nil")
	if ty.Elem != nil {
		return List{Elem: tcx.graphqlTyToTy(ty.Elem)}
	}

	return namedType{Name: ty.NamedType}
}

func (tcx *typechecker) defToTy(def *syn.Definition) Ty {
	switch def.Kind {
	case syn.Scalar:
		switch def.Name {
		case "ID", "String":
			// Technically, ID should always serialize to a string, but not everyone follows the spec.
			return String{}
		case "Int", "Float":
			return Number{}
		case "Boolean":
			return Bool{}
		default:
			// can't any better than returning an Any for a non-builtin scalar
			return Any{}
		}
	case syn.Object:
		ty := Object{Fields: orderedmap.New[string, Ty]()}
		for _, field := range def.Fields {
			// don't convert types recursively as this can lead to infinite recursion
			// just do 1 layer
			ty.Fields.Set(field.Name, tcx.graphqlTyToTy1(field.Type))
		}
		return ty
	case syn.Interface, syn.Union, syn.Enum:
		// TODO
		return Any{}
	case syn.InputObject:
		panic("unreachable, can't query for input types")
	default:
		panic(fmt.Sprintf("unknown definition kind %v", def.Kind))
	}
}

// `expand` at least one layer of the type
func (tcx *typechecker) expand(ty Ty) Ty {
	switch ty := ty.(type) {
	case namedType:
		def, ok := tcx.schema.Types[ty.Name.Value]
		if !ok {
			return tcx.error(ty.Name, fmt.Sprintf("unknown type %v", ty.Name.Value))
		}
		return tcx.defToTy(def)
	default:
		return ty
	}
}

// `flatten` the object type the same way we flatten the response during execution.
func flatten(ty Ty) Ty {
	switch ty := ty.(type) {
	case Object:
		if ty.Fields.Len() == 1 {
			for entry := ty.Fields.Oldest(); entry != nil; entry = entry.Next() {
				return flatten(entry.Value)
			}
		}
		return ty
	case List:
		return List{Elem: flatten(ty.Elem)}
	case Tuple:
		return Tuple{Elems: slice.Map(ty.Elems, flatten)}
	default:
		return ty
	}
}
