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
	scope     *scope
	fragments map[string]*syn.FragmentDefinition
	settings  Settings
}

func (tcx *typechecker) PushScope() {
	tcx.scope = &scope{parent: tcx.scope, bindings: make(map[string]scopeEntry)}
}

func (tcx *typechecker) PopScope() {
	tcx.scope = tcx.scope.parent
}

type scope struct {
	parent   *scope
	bindings map[string]scopeEntry
}

func (s *scope) Lookup(name string) (scopeEntry, bool) {
	if entry, ok := s.bindings[name]; ok {
		return entry, true
	}

	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return scopeEntry{}, false
}

func (s *scope) Bind(name string, ty Ty, pat *syn.NamePat) {
	// we overwrite the name in scope if it already exists (i.e. shadowing is allowed)
	s.bindings[name] = scopeEntry{Ty: ty, Pat: pat}
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
		scope:     &scope{bindings: make(map[string]scopeEntry)},
		fragments: make(map[string]*syn.FragmentDefinition),
		info: Info{
			ExprTypes:         make(map[syn.Expr]Ty),
			BindingTypes:      make(map[*syn.NamePat]Ty),
			NameResolutions:   make(map[*syn.NameExpr]*syn.NamePat),
			VarPatResolutions: make(map[*syn.VarPat]*syn.NamePat),
		},
	}
}

type Info struct {
	Ast          syn.File
	ExprTypes    map[syn.Expr]Ty
	BindingTypes map[*syn.NamePat]Ty
	// NameResolutions maps name expressions to the binding that it references
	NameResolutions map[*syn.NameExpr]*syn.NamePat
	// VarPatResolutions maps $var patterns to the binding that it references
	VarPatResolutions map[*syn.VarPat]*syn.NamePat
	Warnings          Errors
	Errors            Errors
}

func (tcx *typechecker) Check(ast syn.File) Info {
	tcx.info.Ast = ast
	for _, stmt := range ast.Stmts {
		tcx.stmt(stmt)
	}

	tcx.populateWarnings()

	return tcx.info
}

func (tcx *typechecker) populateWarnings() {
	for pat := range tcx.unusedBindings() {
		tcx.warn(pat.Pos(), fmt.Sprintf("unused variable '%s'", pat.Name.Value))
	}
}

func (tcx *typechecker) unusedBindings() map[*syn.NamePat]struct{} {
	// start with all bindings in scope
	unusedBindings := make(map[*syn.NamePat]struct{})
	for pat := range tcx.info.BindingTypes {
		unusedBindings[pat] = struct{}{}
	}

	// remove all bindings that are referenced by name expressions
	for _, binding := range tcx.info.NameResolutions {
		delete(unusedBindings, binding)
	}

	// remove all bindings that are referenced by $var patterns
	for _, binding := range tcx.info.VarPatResolutions {
		delete(unusedBindings, binding)
	}

	return unusedBindings
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

func (tcx *typechecker) warn(pos ast.Position, msg string) {
	tcx.info.Warnings = append(tcx.info.Warnings, Error{pos.Pos(), msg})
}
