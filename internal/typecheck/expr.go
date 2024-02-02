package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/lex"
	"github.com/andyyu2004/gqlt/internal/slice"
	"github.com/andyyu2004/gqlt/syn"
)

func (tcx *typechecker) expr(expr syn.Expr) Ty {
	switch expr := expr.(type) {
	case *syn.NameExpr:
		return Any{}
	case *syn.BinaryExpr:
		return tcx.binaryExpr(expr)
	case *syn.CallExpr:
		return Any{}
	case *syn.IndexExpr:
		return Any{}
	case *syn.LiteralExpr:
		return tcx.literalExpr(expr, expr.Value)
	case *syn.ListExpr:
		return tcx.listExpr(expr)
	case *syn.MatchesExpr:
		return Any{}
	case *syn.ObjectExpr:
		return Any{}
	case *syn.QueryExpr:
		return tcx.queryExpr(expr)
	case *syn.TryExpr:
		return tcx.tryExpr(expr)
	case *syn.UnaryExpr:
		return Any{}
	default:
		panic(fmt.Sprintf("missing case typechecker.expr %T", expr))
	}
}

func (tcx *typechecker) queryExpr(expr *syn.QueryExpr) Ty {
	// todo
	return Any{}
}

func (tcx *typechecker) tryExpr(expr *syn.TryExpr) Ty {
	return Object{
		Fields: map[string]Ty{
			"data": tcx.expr(expr.Expr),
			"errors": List{Elem: Object{
				Fields: map[string]Ty{
					"message": String{},
					"path":    List{Elem: Any{}},
				},
			}},
		},
	}
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
			return List{Elem: ty}
		}

		if !compat(ty, expected) {
			// if the list is not of a uniform type, we assign it a tuple type instead
			return Tuple{Elems: elems}
		}
	}

	return List{Elem: expected}
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
			return tcx.error(expr.Pos(), fmt.Sprintf("cannot compare '%v' to '%v'", lhs, rhs))
		}
		return Bool{}
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
