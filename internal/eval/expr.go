package eval

import (
	"context"
	"fmt"
	"maps"

	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/syn"
)

// keep in sync with typecheck/expr.go
func (e *Executor) eval(ctx context.Context, ecx *executionContext, expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.TryExpr:
		// keep in sync with typecheck/expr.go (tryExpr)
		const dataKey = "data"
		const errorsKey = "errors"
		data, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			if err, ok := err.(catchable); ok {
				return map[string]any{
					dataKey:   data,
					errorsKey: err.catch(),
				}, nil
			}

			return nil, err
		}

		return map[string]any{
			dataKey:   data,
			errorsKey: nil,
		}, nil

	case *syn.QueryExpr:
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
				return nil, errorf(expr.Expr, "cannot negate %T", lhs)
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

		// keep in sync with typecheck/expr.go (binaryExpr)
		switch expr.Op.Kind {
		case lex.Equals2:
			return eq(lhs, rhs), nil
		case lex.BangEqual:
			return !eq(lhs, rhs), nil
		case lex.Plus:
			return add(expr, lhs, rhs)
		case lex.Minus:
			return sub(expr, lhs, rhs)
		case lex.Star:
			return mul(lhs, rhs)
		case lex.Slash:
			return div(lhs, rhs)
		case lex.AngleL:
			return lt(lhs, rhs)
		case lex.AngleLEqual:
			return lte(lhs, rhs)
		case lex.AngleR:
			return gt(lhs, rhs)
		case lex.AngleREqual:
			return gte(lhs, rhs)
		case lex.EqualsTilde:
			return regexMatch(expr, lhs, rhs)
		case lex.BangTilde:
			match, ok := regexMatch(expr, lhs, rhs)
			return !match, ok
		default:
			panic(fmt.Sprintf("missing binary expr eval case: %s", expr.Op))
		}

	case *syn.NameExpr:
		val, ok := ecx.scope.Lookup(expr.Name.Value)
		if !ok {
			return nil, errorf(expr, "reference to undefined variable: %s", expr.Name.Value)
		}

		// copy in order to avoid mutation of the stored values
		switch val := val.(type) {
		case map[string]any:
			return maps.Clone(val), nil
		case []any:
			dst := make([]any, len(val))
			copy(dst, val)
			return dst, nil
		default:
			return val, nil
		}

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
		if base := expr.Base; base != nil {
			val, err := e.eval(ctx, ecx, base)
			if err != nil {
				return nil, err
			}
			var ok bool
			fields, ok = val.(map[string]any)
			if !ok {
				return nil, errorf(base, "object base must be an object, got %T", val)
			}
		}

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

		ecx.PushScope()
		defer ecx.PopScope()
		matches := bindPat(ecx.scope, expr.Pat, val) == nil
		if !matches {
			return false, nil
		}

		if expr.Cond != nil {
			cond, err := e.eval(ctx, ecx, expr.Cond)
			if err != nil {
				return nil, err
			}

			if !truthy(cond) {
				return false, nil
			}
		}

		return true, nil

	case *syn.FieldExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		switch val := val.(type) {
		case map[string]any:
			if val, ok := val[expr.Field.Value]; ok {
				return val, nil
			}
			return nil, nil
		default:
			return nil, errorf(expr.Field, "cannot access field %s on %T", expr.Field.Value, val)
		}

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

		return nil, errorf(expr.Index, "cannot index %T with %T", val, idx)

	case *syn.CallExpr:
		fn, err := e.eval(ctx, ecx, expr.Fn)
		if err != nil {
			return nil, err
		}

		f, ok := fn.(function)
		if !ok {
			return nil, errorf(expr.Fn, "expected function, found %T", fn)
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
