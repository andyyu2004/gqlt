package typecheck

import "github.com/andyyu2004/gqlt/syn"

func (tcx *typechecker) let(stmt *syn.LetStmt) {
	ty := tcx.expr(stmt.Expr)
	tcx.bind(stmt.Pat, ty)
}
