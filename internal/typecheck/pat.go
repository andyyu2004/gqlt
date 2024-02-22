package typecheck

import (
	"fmt"

	"github.com/movio/gqlt/syn"
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
	// This can actually be hit if we encounter an any type somewhere
	// See 66abd8b2a896a9ed3ba6a50439e4639d789585e1
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
		if len(list.Pats) > len(ty.Elems) {
			tcx.error(list.Pos(), fmt.Sprintf("cannot bind tuple with %d elements to a pattern with %d elements", len(list.Pats), len(ty.Elems)))
			return
		}
		for i, pat := range list.Pats {
			if rest, ok := pat.(*syn.RestPat); ok {
				tcx.bind(rest.Pat, Tuple{Elems: ty.Elems[i:]})
			} else {
				tcx.bind(pat, ty.Elems[i])
			}
		}
	case Any, errTy:
		for _, pat := range list.Pats {
			if rest, ok := pat.(*syn.RestPat); ok {
				tcx.bind(rest.Pat, Any{})
				continue
			}
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
				if _, ok := ty.Fields.Get(entry.Key.Value); ok {
					// if it's in the original type, but not in the one we're removing fields from
					tcx.error(entry.Key.Pos(), fmt.Sprintf("field '%s' specified twice", entry.Key.Value))
				} else {
					tcx.error(entry.Key.Pos(), fmt.Sprintf("field '%s' doesn't exist in object type '%s'", entry.Key.Value, ty))
				}
				continue
			}

			tcx.bind(entry.Value, fieldTy)
			// delete the field so we can use the remaining fields to bind the rest pattern to
			_, ok = fieldTys.Delete(entry.Key.Value)
			if !ok {
				panic("we just checked, it should be there..")
			}
		}
	case Any, errTy:
		for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
			tcx.bind(entry.Value, Any{})
		}
	default:
		tcx.error(pat.Pos(), fmt.Sprintf("cannot bind %s to an object pattern", ty))
	}
}

func (tcx *typechecker) bindName(pat *syn.NamePat, ty Ty) {
	tcx.scope.Bind(pat.Name.Value, ty, pat)

	// the name pat itself should be unique (even if the name is shared with other name pats)
	if _, ok := tcx.info.BindingTypes[pat]; ok {
		panic("pattern is being bound twice?")
	}
	tcx.info.BindingTypes[pat] = ty
}
