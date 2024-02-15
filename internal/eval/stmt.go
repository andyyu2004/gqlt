package eval

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
		val, err := e.eval(ctx, ecx, stmt.Expr)
		if err != nil {
			return err
		}

		if !truthy(val) {
			var fmt strings.Builder
			stmt.Expr.Format(&fmt)
			return errorf(stmt, "assertion failed: %v", fmt.String())
		}

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

func (e *Executor) fragment(ctx context.Context, ecx *executionContext, stmt *syn.FragmentStmt) error {
	if _, ok := ecx.scope.fragments[stmt.Fragment.Name.Value]; ok {
		return errorf(stmt, "fragment %s already defined", stmt.Fragment.Name.Value)
	}

	ecx.scope.fragments[stmt.Fragment.Name.Value] = stmt.RawFragment
	return nil
}
