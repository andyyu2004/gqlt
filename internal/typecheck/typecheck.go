package typecheck

import (
	"fmt"
	"strings"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/syn"
)

type Errors []Error

func (errs Errors) Error() string {
	s := new(strings.Builder)
	for _, err := range errs {
		fmt.Fprintf(s, "%v: %s\n", err.Position, err.Msg)
	}
	return s.String()
}

type Error struct {
	Position ast.Position
	Msg      string
}

func (e Error) Message() string {
	return e.Msg
}

func (e Error) Pos() ast.Position {
	return e.Position
}

type typechecker struct {
	schema    *syn.Schema
	info      Info
	scope     map[string]scopeEntry
	fragments map[string]*syn.FragmentDefinition
	settings  Settings
}

type Settings interface {
	Namespace() []string
	Set(key string, value any) error
}

type scopeEntry struct {
	Ty  Ty
	Pat *syn.NamePat
}

// Create a new typechecker.
// Pass the schema to the typechecker to resolve query types against.
// The schema may be nil, in which case the typechecker will typecheck all queries/mutations as any.
func New(schema *syn.Schema, settings Settings) *typechecker {
	return &typechecker{
		schema:    schema,
		settings:  settings,
		scope:     make(map[string]scopeEntry),
		fragments: make(map[string]*syn.FragmentDefinition),
		info: Info{
			ExprTypes:    make(map[syn.Expr]Ty),
			BindingTypes: make(map[*syn.NamePat]Ty),
			Resolutions:  make(map[*syn.NameExpr]*syn.NamePat),
		},
	}
}

type Info struct {
	Ast          syn.File
	ExprTypes    map[syn.Expr]Ty
	BindingTypes map[*syn.NamePat]Ty
	// Resolutions maps name expressions to the binding that it references
	Resolutions map[*syn.NameExpr]*syn.NamePat
	Warnings    Errors
	Errors      Errors
}

func (tcx *typechecker) Check(ast syn.File) Info {
	tcx.info.Ast = ast
	for _, stmt := range ast.Stmts {
		tcx.stmt(stmt)
	}

	return tcx.info
}

// Type error handling invariants:
// - If you construct an `errTy`, you must output a message.
// - If you receive an `errTy`, just handle/propogate it appropriately without emitting any further errors.
// This ensures that we don't emit multiple errors for the same expression, and we don't return an errTy
// without emitting an error message.
func (tcx *typechecker) error(pos ast.HasPosition, msg string) errTy {
	tcx.info.Errors = append(tcx.info.Errors, Error{pos.Pos(), msg})
	return errTy{}
}

// func (tcx *typechecker) warn(pos ast.Position, msg string) {
// 	tcx.info.Warnings = append(tcx.info.Warnings, Error{pos.Pos(), msg})
// }
