package gqlt

import (
	"context"
	"fmt"

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
	vars map[string]any
}

func (e *Executor) Run(ctx context.Context, file syn.File) error {
	ecx := executionContext{vars: map[string]any{}}
	for _, stmt := range file.Stmts {
		switch stmt := stmt.(type) {
		case *syn.LetStmt:
			if err := e.let(ctx, &ecx, stmt); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("missing stmt case: %T", stmt))
		}
	}

	return nil
}

func (e *Executor) let(ctx context.Context, ecx *executionContext, let *syn.LetStmt) error {
	val, err := e.eval(ctx, ecx, let.Expr)
	if err != nil {
		return err
	}

	if err := e.bindPat(ecx, let.Pat, val); err != nil {
		return err
	}

	return nil
}

func (e *Executor) bindPat(ecx *executionContext, pat syn.Pat, val any) error {
	switch pat := pat.(type) {
	case *syn.NamePat:
		ecx.vars[pat.Name] = val
		return nil
	default:
		panic(fmt.Sprintf("missing pattern case: %T", pat))
	}
}

func (e *Executor) eval(ctx context.Context, ecx *executionContext, expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.OperationExpr:
		var data any
		// Pass our local variables directly also as graphql variables
		if err := e.client.Request(ctx, expr.Query, ecx.vars, &data); err != nil {
			return nil, err
		}

		return flatten(data), nil
	default:
		panic(fmt.Sprintf("missing expr case: %T", expr))
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
