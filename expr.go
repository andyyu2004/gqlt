package gqlt

import (
	"context"
	"fmt"

	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/syn"
)

func (e *Executor) eval(ctx context.Context, ecx *executionContext, expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.OperationExpr:
		return e.query(ctx, ecx, expr)

	case *syn.LiteralExpr:
		return expr.Value, nil

	case *syn.UnaryExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		switch expr.Op.Kind {
		case lex.Minus:
			switch lhs := val.(type) {
			case float64:
				return -lhs, nil
			default:
				return nil, fmt.Errorf("cannot negate %T", lhs)
			}
		case lex.Bang, lex.Not:
			return !truthy(val), nil
		default:
			panic(fmt.Sprintf("missing unary expr eval case: %s", expr.Op))
		}

	case *syn.BinaryExpr:
		lhs, err := e.eval(ctx, ecx, expr.Left)
		if err != nil {
			return nil, err
		}

		rhs, err := e.eval(ctx, ecx, expr.Right)
		if err != nil {
			return nil, err
		}

		switch expr.Op.Kind {
		case lex.Equals2:
			return eq(lhs, rhs), nil
		case lex.BangEqual:
			return !eq(lhs, rhs), nil
		case lex.Plus:
			return add(lhs, rhs)
		case lex.Minus:
			return sub(lhs, rhs)
		case lex.Star:
			return mul(lhs, rhs)
		case lex.Slash:
			return div(lhs, rhs)
		default:
			panic(fmt.Sprintf("missing binary expr eval case: %s", expr.Op))
		}

	case *syn.NameExpr:
		val, ok := ecx.scope.Lookup(expr.Name.Value)
		if !ok {
			return nil, fmt.Errorf("reference to undefined variable: %s", expr.Name)
		}

		return val, nil

	case *syn.ListExpr:
		vals := make([]any, len(expr.Exprs))
		for i, expr := range expr.Exprs {
			val, err := e.eval(ctx, ecx, expr)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}

		return vals, nil

	case *syn.ObjectExpr:
		fields := make(map[string]any, expr.Fields.Len())
		for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
			name := entry.Key
			val, err := e.eval(ctx, ecx, entry.Value)
			if err != nil {
				return nil, err
			}
			fields[name.Value] = val
		}

		return fields, nil

	case *syn.MatchesExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		return bindPat(dummyBinder{}, expr.Pat, val) == nil, nil

	case *syn.IndexExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		idx, err := e.eval(ctx, ecx, expr.Index)
		if err != nil {
			return nil, err
		}

		switch val := val.(type) {
		case []any:
			switch idx := idx.(type) {
			case float64:
				if idx < 0 || idx >= float64(len(val)) {
					return nil, nil
				}

				return val[int(idx)], nil
			}
		case map[string]any:
			switch idx := idx.(type) {
			case string:
				if val, ok := val[idx]; ok {
					return val, nil
				}
				return nil, nil
			}
		}

		return nil, fmt.Errorf("cannot index %T with %T", val, idx)

	case *syn.CallExpr:
		fn, err := e.eval(ctx, ecx, expr.Fn)
		if err != nil {
			return nil, err
		}

		f, ok := fn.(function)
		if !ok {
			return nil, fmt.Errorf("expected function, found %T", fn)
		}

		args := make([]any, len(expr.Args))
		for i, arg := range expr.Args {
			arg, err := e.eval(ctx, ecx, arg)
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}

		return f(args)

	default:
		panic(fmt.Sprintf("missing expr eval case: %T", expr))
	}
}
