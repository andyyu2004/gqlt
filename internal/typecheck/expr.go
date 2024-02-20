package typecheck

import (
	"fmt"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/internal/slice"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/syn"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func (tcx *typechecker) expr(expr syn.Expr) Ty {
	ty := func() Ty {
		switch expr := expr.(type) {
		case *syn.NameExpr:
			return tcx.nameExpr(expr)
		case *syn.BinaryExpr:
			return tcx.binaryExpr(expr)
		case *syn.CallExpr:
			return tcx.callExpr(expr)
		case *syn.IndexExpr:
			return tcx.indexExpr(expr)
		case *syn.LiteralExpr:
			return tcx.literalExpr(expr, expr.Value)
		case *syn.ListExpr:
			return tcx.listExpr(expr)
		case *syn.MatchesExpr:
			return tcx.matchesExpr(expr)
		case *syn.ObjectExpr:
			return tcx.objectExpr(expr)
		case *syn.QueryExpr:
			return tcx.queryExpr(expr)
		case *syn.TryExpr:
			return tcx.tryExpr(expr)
		case *syn.UnaryExpr:
			return tcx.unaryExpr(expr)
		case *syn.FieldExpr:
			return tcx.fieldExpr(expr)
		default:
			panic(fmt.Sprintf("missing case typechecker.expr %T", expr))
		}
	}()
	tcx.info.ExprTypes[expr] = ty
	if lib.IsNil(ty) {
		panic(fmt.Sprintf("oops, got nil type when typechecking expr: %v", expr))
	}
	return ty
}

func (tcx *typechecker) unaryExpr(expr *syn.UnaryExpr) Ty {
	ty := tcx.expr(expr.Expr)
	if isErr(ty) {
		return ty
	}

	switch expr.Op.Kind {
	case lex.Bang, lex.Not:
		return Bool{}
	case lex.Minus:
		if _, ok := ty.(Number); ok {
			return Number{}
		}
	}

	return tcx.error(expr.Pos(), fmt.Sprintf("cannot apply '%s' to '%v'", expr.Op.Kind.String(), ty))
}

func (tcx *typechecker) tryExpr(expr *syn.TryExpr) Ty {
	errorFields := orderedmap.New[string, Ty](2)
	errorFields.Set("message", String{})
	errorFields.Set("path", List{Elem: Any{}})

	fields := orderedmap.New[string, Ty](2)
	fields.Set("data", tcx.expr(expr.Expr))
	fields.Set("errors", List{Elem: Object{Fields: errorFields}})

	return Object{Fields: fields}
}

func (tcx *typechecker) nameExpr(expr *syn.NameExpr) Ty {
	if entry, ok := tcx.scope[expr.Name.Value]; ok {
		tcx.info.Resolutions[expr] = entry.Pat
		return entry.Ty
	}
	return tcx.error(expr.Pos(), fmt.Sprintf("unbound name '%s'", expr.Name.Value))
}

func (tcx *typechecker) literalExpr(pos ast.HasPosition, val any) Ty {
	switch val := val.(type) {
	case bool:
		return Bool{}
	case float64:
		return Number{}
	case string:
		return String{}
	case nil:
		return Any{}
	default:
		panic(fmt.Sprintf("missing case typechecker.lit %T", val))
	}
}

func (tcx *typechecker) objectExpr(expr *syn.ObjectExpr) Ty {
	var ty Ty = Object{Fields: orderedmap.New[string, Ty]()}
	if expr.Base != nil {
		ty = tcx.expr(expr.Base)
		if isErr(ty) {
			return ty
		}
	}

	switch ty := ty.(type) {
	case Any, errTy:
		return ty
	case Object:
		for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
			fieldTy := tcx.expr(entry.Value)
			if isErr(ty) {
				return ty
			}

			// overwrite any existing field from the base regardless of prior type
			ty.Fields.Set(entry.Key.Value, fieldTy)
		}

		return ty

	default:
		return tcx.error(expr.Pos(), fmt.Sprintf("object base must be an object, got '%v'", ty))
	}
}

func (tcx *typechecker) listExpr(expr *syn.ListExpr) Ty {
	if len(expr.Exprs) == 0 {
		return Tuple{}
	}

	elems := slice.Map(expr.Exprs, func(e syn.Expr) Ty { return tcx.expr(e) })
	expected := elems[0]

	for _, ty := range elems {
		if isErr(ty) {
			return ty
		}

		if isAny(ty) {
			// could be smarter about this
			// this to avoid a list like `[Bool, Any]` to get assigned type `Bool[]`
			return List{Elem: ty}
		}

		if !compat(ty, expected) {
			// if the list is not of a uniform type, we assign it a tuple type instead
			return Tuple{Elems: elems}
		}
	}

	return List{Elem: expected}
}

func (tcx *typechecker) callExpr(expr *syn.CallExpr) Ty {
	// There is no way of getting the types of functions as they are only available at runtime.
	// tcx.expr(expr.Fn)
	for _, arg := range expr.Args {
		tcx.expr(arg)
	}
	// We just return Any for now as the return type
	return Any{}
}

func (tcx *typechecker) indexExpr(expr *syn.IndexExpr) Ty {
	ty := tcx.expr(expr.Expr)
	if isErr(ty) {
		return ty
	}

	indexTy := tcx.expr(expr.Index)
	if isErr(indexTy) {
		return indexTy
	}

	switch ty := ty.(type) {
	case Object:
		if lit, ok := expr.Index.(*syn.LiteralExpr); ok {
			if k, ok := lit.Value.(string); ok {
				v, ok := ty.Fields.Get(k)
				if !ok {
					return tcx.error(expr.Pos(), fmt.Sprintf("cannot access field '%s' on object %s", k, ty))
				}

				return v
			}
		}

		return Any{}
	case List:
		switch indexTy.(type) {
		case Number:
			return ty.Elem
		default:
			return tcx.error(expr.Pos(), fmt.Sprintf("cannot index '%v' with '%v'", ty, expr.Index))
		}
	case Tuple:
		if lit, ok := expr.Index.(*syn.LiteralExpr); ok {
			if f, ok := lit.Value.(float64); ok {
				i := int(f)
				if i >= 0 && i < len(ty.Elems) {
					return ty.Elems[i]
				}
				return tcx.error(expr.Pos(), fmt.Sprintf("tuple index out of range: %d", i))
			}
		}

		return Any{}
	}

	return tcx.error(expr.Pos(), fmt.Sprintf("cannot index '%v' with '%v'", ty, expr.Index))
}

func (tcx *typechecker) matchesExpr(expr *syn.MatchesExpr) Ty {
	ty := tcx.expr(expr.Expr)
	tcx.bind(expr.Pat, ty)
	return Bool{}
}

func (tcx *typechecker) fieldExpr(expr *syn.FieldExpr) Ty {
	ty := tcx.expr(expr.Expr)
	switch ty := ty.(type) {
	case Object:
		if v, ok := ty.Fields.Get(expr.Field.Value); ok {
			return v
		}
	}

	return tcx.error(expr.Pos(), fmt.Sprintf("cannot access field '%s' on type '%v'", expr.Field.Value, ty))
}

func (tcx *typechecker) binaryExpr(expr *syn.BinaryExpr) Ty {
	lhs := tcx.expr(expr.Left)
	rhs := tcx.expr(expr.Right)

	if isErr(lhs) && isErr(rhs) {
		return errTy{}
	}

	if isAny(lhs) || isAny(rhs) {
		return Any{}
	}

	// keep in sync with eval/expr.go
	switch expr.Op.Kind {
	case lex.Equals2, lex.BangEqual:
		if !compat(lhs, rhs) {
			tcx.error(expr.Pos(), fmt.Sprintf("cannot equate '%v' to '%v'", lhs, rhs))
		}
		return Bool{}
	case lex.AngleLEqual, lex.AngleREqual, lex.AngleL, lex.AngleR:
		switch lhs.(type) {
		case Number, String:
			if !compat(lhs, rhs) {
				return tcx.error(expr.Pos(), fmt.Sprintf("cannot order '%v' against '%v'", lhs, rhs))
			}
		default:
			return tcx.error(expr.Pos(), fmt.Sprintf("cannot order '%v' against '%v'", lhs, rhs))
		}
		return Bool{}

	case lex.BangTilde, lex.EqualsTilde:
		if _, ok := lhs.(String); ok {
			if _, ok := rhs.(String); ok {
				return Bool{}
			}
		}
		return tcx.error(expr.Pos(), fmt.Sprintf("cannot apply '%s' to '%v' and '%v' (expected string and string)", expr.Op.Kind.String(), lhs, rhs))
	case lex.Plus:
		switch lhs := lhs.(type) {
		case Number:
			switch rhs.(type) {
			case Number:
				return Number{}
			}
		case String:
			switch rhs.(type) {
			case String:
				return String{}
			}
		case List:
			switch rhs := rhs.(type) {
			case List:
				if !compat(lhs.Elem, rhs.Elem) {
					return tcx.error(expr.Pos(), fmt.Sprintf("cannot append '%v' to '%v'", lhs, rhs))
				}
				return List{Elem: lhs.Elem}
			}
		}
	case lex.Minus:
		switch lhs.(type) {
		case Number:
			switch rhs.(type) {
			case Number:
				return Number{}
			}
		}
	case lex.Star:
		switch lhs := lhs.(type) {
		case Number:
			switch rhs.(type) {
			case Number:
				return Number{}
			}
		case Tuple:
			if rhs, ok := expr.Right.(*syn.LiteralExpr); ok {
				if f, ok := rhs.Value.(float64); ok {
					n := int(f)
					newType := Tuple{Elems: make([]Ty, 0, len(lhs.Elems)*n)}
					for i := 0; i < n; i++ {
						newType.Elems = append(newType.Elems, lhs.Elems...)
					}
					return newType
				}
			}
		case List:
			switch rhs.(type) {
			case Number:
				return List{Elem: lhs.Elem}
			}
		}
	case lex.Slash:
		switch lhs.(type) {
		case Number:
			switch rhs.(type) {
			case Number:
				return Number{}
			}
		}
	default:
		panic(fmt.Sprintf("missing binary expr typecheck case: %s", expr.Op))
	}

	return tcx.error(expr.Pos(), fmt.Sprintf("cannot apply '%s' to '%v' and '%v'", expr.Op.Kind.String(), lhs, rhs))
}
