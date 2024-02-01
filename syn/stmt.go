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

type FragmentStmt struct {
	ast.Position
	// unparsed graphql string
	// useful for pretty printing without formatting
	RawFragment string
	Fragment    *FragmentDefinition
}

var _ Stmt = FragmentStmt{}

func (stmt FragmentStmt) Children() Children {
	return Children{
		stmt.Fragment.FragmentKw,
		stmt.Fragment.Name,
		stmt.Fragment.OnKw,
		stmt.Fragment.TypeCondition,
	}
}

func (FragmentStmt) isStmt() {}
func (FragmentStmt) isNode() {}

func (stmt FragmentStmt) Dump(w io.Writer) {
	_, _ = io.WriteString(w, stmt.RawFragment)
}

type SetStmt struct {
	ast.Position
	SetKw lex.Token
	Key   lex.Token
	Value Expr
}

var _ Stmt = SetStmt{}

func (stmt SetStmt) Children() Children {
	return Children{stmt.SetKw, stmt.Key, stmt.Value}
}

func (SetStmt) isStmt() {}
func (SetStmt) isNode() {}

func (stmt SetStmt) Dump(w io.Writer) {
	_, _ = io.WriteString(w, "set ")
	_, _ = io.WriteString(w, stmt.Key.Value)
	_, _ = io.WriteString(w, " ")
	stmt.Value.Dump(w)
}

type AssertStmt struct {
	ast.Position
	AssertKw lex.Token
	Expr     Expr
}

var _ Stmt = AssertStmt{}

func (stmt AssertStmt) Children() Children {
	return Children{stmt.AssertKw, stmt.Expr}
}

func (AssertStmt) isStmt() {}
func (AssertStmt) isNode() {}

func (stmt AssertStmt) Dump(w io.Writer) {
	_, _ = io.WriteString(w, "assert ")
	stmt.Expr.Dump(w)
}

type LetStmt struct {
	ast.Position
	LetKw  lex.Token
	Pat    Pat
	Equals lex.Token
	Expr   Expr
}

var _ Stmt = LetStmt{}

func (let LetStmt) Children() Children {
	return Children{let.LetKw, let.Pat, let.Equals, let.Expr}
}

func (LetStmt) isStmt() {}
func (LetStmt) isNode() {}

func (let LetStmt) Dump(w io.Writer) {
	_, _ = io.WriteString(w, "let ")
	let.Pat.Dump(w)
	_, _ = io.WriteString(w, " = ")
	let.Expr.Dump(w)
}
