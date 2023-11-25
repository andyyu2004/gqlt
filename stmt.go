package gqlt

import (
	"context"
	"fmt"

	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/syn"
)

func (e *Executor) stmt(ctx context.Context, ecx *executionContext, stmt syn.Stmt) error {
	switch stmt := stmt.(type) {
	case *syn.LetStmt:
		if err := e.let(ctx, ecx, stmt); err != nil {
			return err
		}
	case *syn.ExprStmt:
		if _, err := e.eval(ctx, ecx, stmt.Expr); err != nil {
			return err
		}

	case *syn.SetStmt:
		val, err := e.eval(ctx, ecx, stmt.Value)
		if err != nil {
			return err
		}

		if err := ecx.settings.Set(stmt.Key, val); err != nil {
			return err
		}

	case *syn.AssertStmt:
		bin, ok := stmt.Expr.(*syn.BinaryExpr)
		if ok && bin.Op == lex.Equals2 {
			// special case for common equality assertions to have a better error message
			lhs, err := e.eval(ctx, ecx, bin.Left)
			if err != nil {
				return err
			}

			rhs, err := e.eval(ctx, ecx, bin.Right)
			if err != nil {
				return err
			}

			if !eq(lhs, rhs) {
				return fmt.Errorf("assertion failed: %v != %v", lhs, rhs)
			}
		} else {
			val, err := e.eval(ctx, ecx, stmt.Expr)
			if err != nil {
				return err
			}

			if !truthy(val) {
				return fmt.Errorf("assertion failed")
			}

		}

	default:
		panic(fmt.Sprintf("missing stmt eval case: %T", stmt))
	}

	return nil
}

func (e *Executor) let(ctx context.Context, ecx *executionContext, let *syn.LetStmt) error {
	val, err := e.eval(ctx, ecx, let.Expr)
	if err != nil {
		return err
	}

	if err := bindPat(ecx.scope, let.Pat, val); err != nil {
		return err
	}

	return nil
}
