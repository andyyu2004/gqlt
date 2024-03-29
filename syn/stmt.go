package syn

import (
	"io"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/lex"
)

type Stmt interface {
	Node
	isStmt()
}

type ExprStmt struct {
	Expr Expr
}

var _ Stmt = ExprStmt{}

func (stmt ExprStmt) Pos() ast.Position {
	return stmt.Expr.Pos()
}

func (stmt ExprStmt) Children() Children {
	return Children{stmt.Expr}
}

func (ExprStmt) isStmt() {}
func (ExprStmt) isNode() {}

func (stmt ExprStmt) Format(w io.Writer) {
	stmt.Expr.Format(w)
}

type FragmentStmt struct {
	ast.Position
	// unparsed graphql string
	// useful for pretty printing without formatting
	RawFragment string
	Definition  *FragmentDefinition
}

var _ Stmt = FragmentStmt{}

func (stmt FragmentStmt) Children() Children {
	return Children{stmt.Definition}
}

func (FragmentStmt) isStmt() {}
func (FragmentStmt) isNode() {}

func (stmt FragmentStmt) Format(w io.Writer) {
	_, _ = io.WriteString(w, stmt.RawFragment)
}

type SetStmt struct {
	ast.Position
	SetKw    lex.Token
	Variable lex.Token
	Expr     Expr
}

var _ Stmt = SetStmt{}

func (stmt SetStmt) Children() Children {
	return Children{stmt.SetKw, stmt.Variable, stmt.Expr}
}

func (SetStmt) isStmt() {}
func (SetStmt) isNode() {}

func (stmt SetStmt) Format(w io.Writer) {
	_, _ = io.WriteString(w, "set ")
	_, _ = io.WriteString(w, stmt.Variable.Value)
	_, _ = io.WriteString(w, " ")
	stmt.Expr.Format(w)
}

type AssertStmt struct {
	AssertKw lex.Token
	Expr     Expr
}

var _ Stmt = AssertStmt{}

func (stmt AssertStmt) Pos() ast.Position {
	return stmt.AssertKw.Merge(stmt.Expr.Pos())
}

func (stmt AssertStmt) Children() Children {
	return Children{stmt.AssertKw, stmt.Expr}
}

func (AssertStmt) isStmt() {}
func (AssertStmt) isNode() {}

func (stmt AssertStmt) Format(w io.Writer) {
	_, _ = io.WriteString(w, "assert ")
	stmt.Expr.Format(w)
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

func (let LetStmt) Format(w io.Writer) {
	_, _ = io.WriteString(w, "let ")
	let.Pat.Format(w)
	_, _ = io.WriteString(w, " = ")
	let.Expr.Format(w)
}

type UseStmt struct {
	UseKw lex.Token
	Path  lex.Token
}

var _ Stmt = UseStmt{}

func (stmt UseStmt) Pos() ast.Position {
	return stmt.UseKw.Merge(stmt.Path.Position)
}

func (stmt UseStmt) Children() Children {
	return Children{stmt.UseKw, stmt.Path}
}

func (UseStmt) isStmt() {}
func (UseStmt) isNode() {}

func (stmt UseStmt) Format(w io.Writer) {
	_, _ = io.WriteString(w, "use ")
	_, _ = io.WriteString(w, stmt.Path.Value)
}
