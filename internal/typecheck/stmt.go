package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/syn"
)

func (tcx *typechecker) stmt(stmt syn.Stmt) {
	switch stmt := stmt.(type) {
	case *syn.ExprStmt:
		tcx.expr(stmt.Expr)
	case *syn.LetStmt:
		tcx.let(stmt)
	case *syn.AssertStmt:
		ty := tcx.expr(stmt.Expr)
		_ = ty
		// allow any type to be asserted against for now
	case *syn.SetStmt:
		tcx.set(stmt)
	case *syn.FragmentStmt:
		tcx.fragment(stmt)
	case *syn.UseStmt:
	default:
		panic(fmt.Sprintf("missing case typechecker.stmt %T", stmt))
	}
}

func (tcx *typechecker) fragment(stmt *syn.FragmentStmt) {
	if _, ok := tcx.fragments[stmt.Fragment.Name.Value]; ok {
		tcx.error(stmt, fmt.Sprintf("fragment %q already defined", stmt.Fragment.Name.Value))
		return
	}
	tcx.fragments[stmt.Fragment.Name.Value] = stmt.Fragment
}

func (tcx *typechecker) let(stmt *syn.LetStmt) {
	ty := tcx.expr(stmt.Expr)
	tcx.bind(stmt.Pat, ty)
}

func (tcx *typechecker) set(stmt *syn.SetStmt) {
	ty := tcx.expr(stmt.Expr)
	// keep in sync with eval
	switch stmt.Variable.Value {
	case "namespace":
		// HACK: We actually need to evaluate and set this one because it affects how we typecheck queries...
		switch ty := ty.(type) {
		case String:
			switch expr := stmt.Expr.(type) {
			case *syn.LiteralExpr:
				if err := tcx.settings.Set("namespace", expr.Value); err != nil {
					tcx.error(stmt, fmt.Sprintf("failed to set namespace: %v", err))
				}
			}
		case List:
			if _, ok := ty.Elem.(String); ok {
				v, err := constEval(stmt.Expr)
				if err != nil {
					tcx.error(stmt, fmt.Sprintf("failed to set namespace: %v", err))
				} else if err := tcx.settings.Set("namespace", v); err != nil {
					tcx.error(stmt, fmt.Sprintf("failed to set namespace: %v", err))
				}
				return
			}
			tcx.error(stmt, fmt.Sprintf("expected list of strings as value for %q, found %s", stmt.Variable.Value, ty))
		default:
			tcx.error(stmt, fmt.Sprintf("expected string or list of strings as value for %q, found %s", stmt.Variable.Value, ty))
		}
	default:
		tcx.error(stmt, fmt.Sprintf("unknown set variable %q", stmt.Variable.Value))
	}
}

// "compile time" evaluation, used by the typechecker to evaluate expressions such as literals for setting the namespace
func constEval(expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.LiteralExpr:
		return expr.Value, nil
	case *syn.ListExpr:
		vals := make([]any, len(expr.Exprs))
		for i, expr := range expr.Exprs {
			val, err := constEval(expr)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}

		return vals, nil
	case *syn.ObjectExpr:
		object := make(map[string]any, expr.Fields.Len())
		for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
			name := entry.Key
			val, err := constEval(entry.Value)
			if err != nil {
				return nil, err
			}
			object[name.Value] = val
		}
		return object, nil
	}

	return nil, fmt.Errorf("cannot const eval %T", expr)
}
