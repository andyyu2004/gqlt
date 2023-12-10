package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/lex"
)

type Stmt interface {
	Node
	isStmt()
}

type ExprStmt struct {
	ast.Position
	Expr Expr
}

var _ Stmt = ExprStmt{}

func (stmt ExprStmt) Children() Children {
	return Children{stmt.Expr}
}

func (ExprStmt) isStmt() {}
func (ExprStmt) isNode() {}

func (stmt ExprStmt) Dump(w io.Writer) {
	stmt.Expr.Dump(w)
}

type SetStmt struct {
	ast.Position
	Key   lex.Token
	Value Expr
}

var _ Stmt = SetStmt{}

func (stmt SetStmt) Children() Children {
	return Children{stmt.Key, stmt.Value}
}

func (SetStmt) isStmt() {}
func (SetStmt) isNode() {}

func (stmt SetStmt) Dump(w io.Writer) {
	io.WriteString(w, "set ")
	io.WriteString(w, stmt.Key.Value)
	io.WriteString(w, " ")
	stmt.Value.Dump(w)
}

type AssertStmt struct {
	ast.Position
	Expr Expr
}

var _ Stmt = AssertStmt{}

func (stmt AssertStmt) Children() Children {
	return Children{stmt.Expr}
}

func (AssertStmt) isStmt() {}
func (AssertStmt) isNode() {}

func (stmt AssertStmt) Dump(w io.Writer) {
	io.WriteString(w, "assert ")
	stmt.Expr.Dump(w)
}

type LetStmt struct {
	ast.Position
	Pat  Pat
	Expr Expr
}

var _ Stmt = LetStmt{}

func (let LetStmt) Children() Children {
	return Children{let.Pat, let.Expr}
}

func (LetStmt) isStmt() {}
func (LetStmt) isNode() {}

func (let LetStmt) Dump(w io.Writer) {
	io.WriteString(w, "let ")
	let.Pat.Dump(w)
	io.WriteString(w, " = ")
	let.Expr.Dump(w)
}
