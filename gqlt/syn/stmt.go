package syn

import "io"

type Stmt interface {
	Node
	isStmt()
}

type ExprStmt struct {
	Expr Expr
}

var _ Stmt = ExprStmt{}

func (ExprStmt) isStmt() {}
func (ExprStmt) isNode() {}

func (stmt ExprStmt) Dump(w io.Writer) {
	stmt.Expr.Dump(w)
}

type LetStmt struct {
	Pat  Pat
	Expr Expr
}

var _ Stmt = LetStmt{}

func (LetStmt) isStmt() {}
func (LetStmt) isNode() {}
func (let LetStmt) Dump(w io.Writer) {
	io.WriteString(w, "let ")
	let.Pat.Dump(w)
	io.WriteString(w, " = ")
	let.Expr.Dump(w)
}
