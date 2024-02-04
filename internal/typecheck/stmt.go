package typecheck

import (
	"fmt"

	"github.com/andyyu2004/gqlt/syn"
)

func (tcx *typechecker) let(stmt *syn.LetStmt) {
	ty := tcx.expr(stmt.Expr)
	tcx.bind(stmt.Pat, ty)
}

func (tcx *typechecker) set(stmt *syn.SetStmt) {
	ty := tcx.expr(stmt.Expr)
	// keep in sync with eval
	switch stmt.Variable.Value {
	case "namespace":
		switch ty := ty.(type) {
		case String:
		case List:
			if _, ok := ty.Elem.(String); ok {
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
