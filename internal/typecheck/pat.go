package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/syn"
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
		panic("rest pattern should have special handling in list and object cases")
	default:
	}
}

func (tcx *typechecker) bindList(pat *syn.ListPat, ty Ty) {
	switch ty := ty.(type) {
	case List:
	case Tuple:
	case Any:
	default:
		tcx.error(pat.Pos(), fmt.Sprintf("cannot bind %s to a list pattern", ty))
	}
}

func (tcx *typechecker) bindObject(pat *syn.ObjectPat, ty Ty) {
	switch ty := ty.(type) {
	case Object:
	case Any:
	default:
		tcx.error(pat.Pos(), fmt.Sprintf("cannot bind %s to an object pattern", ty))
	}
}

func (tcx *typechecker) bindName(pat *syn.NamePat, ty Ty) {}
