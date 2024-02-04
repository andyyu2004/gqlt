package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/internal/slice"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/syn"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func (tcx *typechecker) queryExpr(expr *syn.QueryExpr) Ty {
	if tcx.schema == nil {
		// Impossible to typecheck a query without a schema.
		// It would be possible to approximate it except for the fact we have no way of telling when a field is a list or not.
		// Could do it based off whether the field name is plural or not, but that's only a heuristic and disgusting.
		return Any{}
	}

	ty := tcx.inferQueryType(expr.Operation)
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

	return tcx.inferSelectionSetType(def, operation.SelectionSet)
}

func (tcx *typechecker) inferSelectionSetType(def *syn.Definition, selectionSet syn.SelectionSet) Ty {
	outTy := Object{Fields: orderedmap.New[string, Ty]()}
	for _, selection := range selectionSet {
		switch selection := selection.(type) {
		case *syn.Field:
			field := def.Fields.ForName(selection.Name.Value)
			var fieldTy Ty
			if field == nil {
				fieldTy = tcx.error(selection, fmt.Sprintf("field %v does not exist on type %v", selection.Name.Value, def.Name))
			} else {
				fieldTy = tcx.graphqlTyToTy(field.Type)
			}
			outTy.Fields.Set(field.Name, fieldTy)
		case *syn.FragmentSpread:
			return Any{}
		case *syn.InlineFragment:
			return Any{}
		}
	}

	return outTy
}

func (tcx *typechecker) graphqlTyToTy(ty *syn.Type) Ty {
	lib.Assert(ty != nil, "type must not be nil")
	if ty.Elem != nil {
		return List{Elem: tcx.graphqlTyToTy(ty.Elem)}
	}

	def, ok := tcx.schema.Types[ty.NamedType.Value]
	if !ok {
		return tcx.error(ty, fmt.Sprintf("unknown type %v", ty.NamedType.Value))
	}

	return tcx.defToTy(def)
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
			ty.Fields.Set(field.Name, tcx.graphqlTyToTy(field.Type))
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
