package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/syn"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func (tcx *typechecker) bind(pat syn.Pat, ty Ty) {
	switch pat := pat.(type) {
	case *syn.WildcardPat:
	case *syn.LiteralPat:
	case *syn.NamePat:
		tcx.bindName(pat, ty)
	case *syn.ListPat:
		tcx.bindList(pat, ty)
	case *syn.ObjectPat:
		tcx.bindObject(pat, ty)
	case *syn.RestPat:
		// panic("rest pattern should have special handling in list and object cases")
	default:
	}
}

func (tcx *typechecker) bindList(list *syn.ListPat, ty Ty) {
	switch ty := ty.(type) {
	case List:
		for i, pat := range list.Pats {
			if rest, ok := pat.(*syn.RestPat); ok {
				tcx.bind(rest.Pat, ty)
				if i != len(list.Pats)-1 {
					panic("rest pattern must be the last pattern in a list (should be checked by parser)")
				}
			} else {
				tcx.bind(pat, ty.Elem)
			}
		}
	case Tuple:
		if len(list.Pats) != len(ty.Elems) {
			tcx.error(list.Pos(), fmt.Sprintf("expected %d elements, found %d", len(ty.Elems), len(list.Pats)))
			return
		}

		for i, pat := range list.Pats {
			if rest, ok := pat.(*syn.RestPat); ok {
				tcx.bind(rest.Pat, Tuple{Elems: ty.Elems[i:]})
			} else {
				tcx.bind(pat, ty.Elems[i])
			}
		}
	case Any:
		for _, pat := range list.Pats {
			tcx.bind(pat, Any{})
		}
	default:
		tcx.error(list.Pos(), fmt.Sprintf("cannot bind %s to a list pattern", ty))
	}
}

func (tcx *typechecker) bindObject(pat *syn.ObjectPat, ty Ty) {
	switch ty := ty.(type) {
	case Object:
		// copying the fields as we will mutate it
		fieldTys := orderedmap.New[string, Ty](ty.Fields.Len())
		for entry := ty.Fields.Oldest(); entry != nil; entry = entry.Next() {
			fieldTys.Set(entry.Key, entry.Value)
		}

		for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
			if rest, ok := entry.Value.(*syn.RestPat); ok {
				tcx.bind(rest.Pat, Object{Fields: fieldTys})
				continue
			}

			fieldTy, ok := fieldTys.Get(entry.Key.Value)
			if !ok {
				tcx.error(pat.Pos(), fmt.Sprintf("field '%s' not found in object", entry.Key.Value))
				continue
			}

			tcx.bind(entry.Value, fieldTy)
			// delete the field so we can use the remaining fields to bind the rest pattern to
			_, ok = fieldTys.Delete(entry.Key.Value)
			if !ok {
				panic("we just checked, it should be there..")
			}
		}
	case Any:
		for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
			tcx.bind(entry.Value, Any{})
		}
	default:
		tcx.error(pat.Pos(), fmt.Sprintf("cannot bind %s to an object pattern", ty))
	}
}

func (tcx *typechecker) bindName(pat *syn.NamePat, ty Ty) {
	// we overwrite the name in scope if it already exists (i.e. shadowing is allowed)
	tcx.scope[pat.Name.Value] = ty

	// the name pat itself should be unique (even if the name is shared with other name pats)
	if _, ok := tcx.info.BindingTypes[pat]; ok {
		panic("pattern is being bound twice?")
	}
	tcx.info.BindingTypes[pat] = ty
}
