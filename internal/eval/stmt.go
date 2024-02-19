package eval

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/syn"
)

func (e *Executor) stmt(ctx context.Context, ecx *executionContext, stmt syn.Stmt) error {
	switch stmt := stmt.(type) {
	case *syn.UseStmt:
		if err := e.use(ctx, ecx, stmt); err != nil {
			return err
		}

	case *syn.LetStmt:
		if err := e.let(ctx, ecx, stmt); err != nil {
			return err
		}
	case *syn.FragmentStmt:
		if err := e.fragment(ctx, ecx, stmt); err != nil {
			return err
		}
	case *syn.ExprStmt:
		if _, err := e.eval(ctx, ecx, stmt.Expr); err != nil {
			return err
		}

	case *syn.SetStmt:
		val, err := e.eval(ctx, ecx, stmt.Expr)
		if err != nil {
			return err
		}

		if err := ecx.settings.Set(stmt.Variable.Value, val); err != nil {
			return err
		}

	case *syn.AssertStmt:
		return e.assert(ctx, ecx, stmt)

	default:
		panic(fmt.Sprintf("missing stmt eval case: %T", stmt))
	}

	return nil
}

func (e *Executor) use(ctx context.Context, ecx *executionContext, use *syn.UseStmt) error {
	value := use.Path.Value
	// default to .gqlt extension
	if filepath.Ext(value) == "" {
		value += Ext
	}

	// the `use` path is relative to the calling file's directory
	path := filepath.Join(filepath.Dir(ecx.path), value)
	return e.RunFile(ctx, ecx.client, path)
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

func (e *Executor) fragment(_ context.Context, ecx *executionContext, stmt *syn.FragmentStmt) error {
	if _, ok := ecx.scope.fragments[stmt.Definition.Name.Value]; ok {
		return errorf(stmt, "fragment %s already defined", stmt.Definition.Name.Value)
	}

	ecx.scope.fragments[stmt.Definition.Name.Value] = stmt
	return nil
}

func (e *Executor) assert(ctx context.Context, ecx *executionContext, stmt *syn.AssertStmt) error {
	// Try to provide nice assertion failure messages for certain common cases
	switch expr := stmt.Expr.(type) {
	case *syn.MatchesExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return err
		}

		if err := bindPat(dummyBinder{}, expr.Pat, val); err != nil {
			var msg string
			if e, ok := err.(Error); ok {
				// need to avoid reporting positions twice
				msg = e.Message()
			} else {
				msg = err.Error()
			}

			return errorf(stmt, "assertion failed: %v", msg)
		}

		return nil

	case *syn.BinaryExpr:
		if expr.Op.Kind == lex.Equals2 {
			lhs, err := e.eval(ctx, ecx, expr.Left)
			if err != nil {
				return err
			}

			rhs, err := e.eval(ctx, ecx, expr.Right)
			if err != nil {
				return err
			}

			diff := cmp.Diff(lhs, rhs)
			if diff != "" {
				return errorf(stmt, "assertion failed: %v", diff)
			}

			return nil
		}
	}

	val, err := e.eval(ctx, ecx, stmt.Expr)
	if err != nil {
		return err
	}

	if !truthy(val) {
		var fmt strings.Builder
		stmt.Expr.Format(&fmt)
		return errorf(stmt, "assertion failed: %v", fmt.String())
	}

	return nil
}
