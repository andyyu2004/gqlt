package gqlt

import (
	"context"
	"fmt"
	"reflect"

	"andyyu2004/gqlt/lex"
	"andyyu2004/gqlt/syn"
)

type Client interface {
	Request(ctx context.Context, query string, variables map[string]any, out any) error
}

type Executor struct{ client Client }

type Option func(*Executor)

func New(client Client, opts ...Option) *Executor {
	return &Executor{client}
}

type executionContext struct {
	scope *scope
}

type scope struct {
	parent *scope
	vars   map[string]any
}

func (s *scope) bind(name string, val any) {
	s.vars[name] = val
}

func (s *scope) Lookup(name string) (any, bool) {
	val, ok := s.vars[name]
	if !ok && s.parent != nil {
		return s.parent.Lookup(name)
	}

	return val, ok
}

func (s *scope) gqlVars() map[string]any {
	vars := map[string]any{}

	for name, val := range s.vars {
		// We avoid passing in function values as graphql variables.
		// This is only a shallow check so we can still pass in a map containing functions for example.
		switch reflect.ValueOf(val).Kind() {
		case reflect.Func:
			continue
		default:
			vars[name] = val
		}
	}

	if s.parent != nil {
		for name, val := range s.parent.gqlVars() {
			if _, ok := vars[name]; ok {
				// don't overwrite shadowed variables
				continue
			}

			vars[name] = val
		}
	}

	return vars
}

type function func(args []any) (any, error)

func checkArity(arity int, args []any) error {
	if len(args) != arity {
		return fmt.Errorf("expected %d args, found %d", arity, len(args))
	}

	return nil
}

func truthy(val any) bool {
	switch val := val.(type) {
	case nil:
		return false
	case bool:
		return val
	case int:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

var builtinScope = &scope{
	vars: map[string]any{
		"example": function(func(args []any) (any, error) {
			if err := checkArity(1, args); err != nil {
				return nil, err
			}

			if !truthy(args[0]) {
				return nil, fmt.Errorf("assertion failed")
			}

			return nil, nil
		}),
	},
}

func (e *Executor) Run(ctx context.Context, file syn.File) error {
	ecx := &executionContext{
		scope: &scope{
			parent: builtinScope,
			vars:   map[string]any{},
		},
	}

	for _, stmt := range file.Stmts {
		switch stmt := stmt.(type) {
		case *syn.LetStmt:
			if err := e.let(ctx, ecx, stmt); err != nil {
				return err
			}
		case *syn.ExprStmt:
			if _, err := e.eval(ctx, ecx, stmt.Expr); err != nil {
				return err
			}

		case *syn.AssertStmt:
			val, err := e.eval(ctx, ecx, stmt.Expr)
			if err != nil {
				return err
			}

			if !truthy(val) {
				return fmt.Errorf("assertion failed")
			}

		default:
			panic(fmt.Sprintf("missing stmt eval case: %T", stmt))
		}
	}

	return nil
}

func (e *Executor) let(ctx context.Context, ecx *executionContext, let *syn.LetStmt) error {
	val, err := e.eval(ctx, ecx, let.Expr)
	if err != nil {
		return err
	}

	if err := bindPat(ecx.scope, let.Pat, val); err != nil {
		return err
	}

	return nil
}

type binder interface {
	bind(name string, val any)
}

type dummyBinder struct{}

var _ binder = dummyBinder{}

func (dummyBinder) bind(string, any) {}

func bindPat(binder binder, pat syn.Pat, val any) error {
	switch pat := pat.(type) {
	case *syn.NamePat:
		binder.bind(pat.Name, val)
		return nil
	case *syn.ObjectPat:
		vals, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot bind object pattern to value: %T", val)
		}
		return bindObjectPat(binder, pat, vals)
	case *syn.LiteralPat:
		if pat.Value != val {
			return fmt.Errorf("literal pattern does not match value: %v != %v", pat.Value, val)
		}
		return nil
	default:
		panic(fmt.Sprintf("missing pattern bind case: %T", pat))
	}
}

func bindObjectPat(binder binder, pat *syn.ObjectPat, values map[string]any) error {
	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		name := entry.Key
		pat := entry.Value

		val, ok := values[name]
		if !ok {
			return fmt.Errorf("object missing field specified in pattern %s", name)
		}

		if err := bindPat(binder, pat, val); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) eval(ctx context.Context, ecx *executionContext, expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.OperationExpr:
		var data any
		// Pass our local variables directly also as graphql variables
		if err := e.client.Request(ctx, expr.Query, ecx.scope.gqlVars(), &data); err != nil {
			return nil, err
		}

		return flatten(data), nil

	case *syn.LiteralExpr:
		return expr.Value, nil

	case *syn.BinaryExpr:
		lhs, err := e.eval(ctx, ecx, expr.Left)
		if err != nil {
			return nil, err
		}

		rhs, err := e.eval(ctx, ecx, expr.Right)
		if err != nil {
			return nil, err
		}

		switch expr.Op {
		case lex.Equals2:
			return reflect.DeepEqual(lhs, rhs), nil
		case lex.BangEqual:
			return !reflect.DeepEqual(lhs, rhs), nil
		default:
			panic(fmt.Sprintf("missing binary op eval case: %s", expr.Op))
		}

	case *syn.NameExpr:
		val, ok := ecx.scope.Lookup(expr.Name)
		if !ok {
			return nil, fmt.Errorf("reference to undefined variable: %s", expr.Name)
		}

		return val, nil

	case *syn.ObjectExpr:
		fields := make(map[string]any, expr.Fields.Len())
		for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
			name := entry.Key
			val, err := e.eval(ctx, ecx, entry.Value)
			if err != nil {
				return nil, err
			}
			fields[name] = val
		}

		return fields, nil

	case *syn.MatchesExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		return bindPat(dummyBinder{}, expr.Pat, val) == nil, nil

	case *syn.CallExpr:
		fn, err := e.eval(ctx, ecx, expr.Fn)
		if err != nil {
			return nil, err
		}

		f, ok := fn.(function)
		if !ok {
			return nil, fmt.Errorf("expected function, found %T", fn)
		}

		args := make([]any, len(expr.Args))
		for i, arg := range expr.Args {
			arg, err := e.eval(ctx, ecx, arg)
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}

		return f(args)

	default:
		panic(fmt.Sprintf("missing expr eval case: %T", expr))
	}
}

// flatten removes unnecessary nesting in a (hopefully) intuitive way from the graphql response
func flatten(data any) any {
	switch data := data.(type) {
	case map[string]any:
		if len(data) == 1 {
			for _, v := range data {
				return flatten(v)
			}
		}
		return data
	case []any:
		// recursively flatten elements of arrays
		xs := make([]any, len(data))
		for i, v := range data {
			xs[i] = flatten(v)
		}
		return xs
	default:
		return data
	}
}
