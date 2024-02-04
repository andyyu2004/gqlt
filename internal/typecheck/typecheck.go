package typecheck

import (
	"fmt"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/syn"
)

type Errors []Error

func (errs Errors) Error() string {
	s := new(strings.Builder)
	for _, err := range errs {
		fmt.Fprintf(s, "%v: %s\n", err.Pos, err.Msg)
	}
	return s.String()
}

type Error struct {
	Pos ast.Position
	Msg string
}

type typechecker struct {
	info  Info
	scope map[string]scopeEntry
}

type scopeEntry struct {
	Ty  Ty
	Pat *syn.NamePat
}

func New() *typechecker {
	return &typechecker{
		scope: make(map[string]scopeEntry),
		info: Info{
			ExprTypes:    make(map[syn.Expr]Ty),
			BindingTypes: make(map[*syn.NamePat]Ty),
			Resolutions:  make(map[*syn.NameExpr]*syn.NamePat),
		},
	}
}

type Info struct {
	ExprTypes    map[syn.Expr]Ty
	BindingTypes map[*syn.NamePat]Ty
	// Resolutions maps name expressions to the binding that it references
	Resolutions map[*syn.NameExpr]*syn.NamePat
	Warnings    Errors
	Errors      Errors
}

func (tcx *typechecker) Check(ast syn.File) Info {
	for _, stmt := range ast.Stmts {
		tcx.stmt(stmt)
	}

	return tcx.info
}

func (tcx *typechecker) stmt(stmt syn.Stmt) {
	switch stmt := stmt.(type) {
	case *syn.ExprStmt:
		tcx.expr(stmt.Expr)
	case *syn.LetStmt:
		tcx.let(stmt)
	case *syn.AssertStmt:
		tcx.expr(stmt.Expr)
	case *syn.SetStmt:
		tcx.expr(stmt.Expr)
	case *syn.FragmentStmt, *syn.UseStmt:
	default:
		panic(fmt.Sprintf("missing case typechecker.stmt %T", stmt))
	}
}

func (tcx *typechecker) error(pos ast.Position, msg string) errTy {
	tcx.info.Errors = append(tcx.info.Errors, Error{pos.Pos(), msg})
	return errTy{}
}

func (tcx *typechecker) warn(pos ast.Position, msg string) {
	tcx.info.Warnings = append(tcx.info.Warnings, Error{pos.Pos(), msg})
}
