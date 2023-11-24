package gqlt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/formatter"
	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/parser"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/graph-gophers/graphql-go"
)

// A thread-safe graphql client
type Client interface {
	Request(ctx context.Context, req Request, out any) error
}

type Request struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type GraphQLGophersClient struct {
	Schema *graphql.Schema
}

var _ Client = GraphQLGophersClient{}

func (a GraphQLGophersClient) Request(ctx context.Context, req Request, out any) error {
	res := a.Schema.Exec(ctx, req.Query, req.OperationName, req.Variables)
	if len(res.Errors) > 0 {
		errs := make([]error, 0, len(res.Errors))
		for _, err := range res.Errors {
			errs = append(errs, err)
		}

		return errors.Join(errs...)
	}

	return json.Unmarshal([]byte(res.Data), out)
}

type HTTPClient struct {
	Handler http.Handler
	Headers http.Header
	URL     string
}

var _ Client = GraphQLGophersClient{}

func (c HTTPClient) Request(ctx context.Context, req Request, out any) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return fmt.Errorf("unable to encode request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, &buf)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	if c.Headers != nil {
		httpReq.Header = c.Headers.Clone()
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Accept", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, httpReq)

	type Response struct {
		Errors GraphqlErrors `json:"errors"`
		Data   interface{}
	}

	graphqlResponse := Response{Data: out}

	if err := json.NewDecoder(w.Body).Decode(&graphqlResponse); err != nil {
		return err
	}

	if len(graphqlResponse.Errors) > 0 {
		return graphqlResponse.Errors
	}

	return nil
}

type Response struct {
	Errors GraphqlErrors `json:"errors"`
	Data   interface{}
}

type GraphqlErrors []GraphqlError

type GraphqlError struct {
	Message string   `json:"message"`
	Path    ast.Path `json:"path,omitempty"`
}

func (e GraphqlErrors) Error() string {
	var errs []string
	for _, err := range e {
		errs = append(errs, err.Message)
	}
	return strings.Join(errs, ",")
}

type Executor struct{ client Client }

func New(client Client) *Executor {
	return &Executor{client}
}

type settings struct {
	namespace []string
}

func (s *settings) Set(key string, val any) error {
	switch key {
	case "namespace":
		switch val := val.(type) {
		case string:
			s.namespace = strings.Split(val, "/")
			return nil
		case []any:
			parts := make([]string, len(val))
			for i, v := range val {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("expected string elements in namespace list, found %T", v)
				}
				parts[i] = s
			}
			s.namespace = parts
			return nil
		default:
			return fmt.Errorf("expected slash-separated string or list of strings for namespace setting, found %T", val)
		}
	default:
		return fmt.Errorf("unknown setting: %s", key)
	}
}

type executionContext struct {
	scope    *scope
	settings settings
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

func eq(lhs, rhs any) bool {
	return reflect.DeepEqual(lhs, rhs)
}

func add(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs + rhs, nil
		}
	case string:
		switch rhs := rhs.(type) {
		case string:
			return lhs + rhs, nil
		}
	case []any:
		switch rhs := rhs.(type) {
		case []any:
			return append(lhs, rhs...), nil
		}
	}

	return nil, fmt.Errorf("cannot add %T and %T", lhs, rhs)
}

func sub(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs - rhs, nil
		}
	}

	return nil, fmt.Errorf("cannot subtract %T and %T", lhs, rhs)
}

func mul(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs * rhs, nil
		}
	case []any:
		switch rhs := rhs.(type) {
		case float64:
			n := int(rhs)
			copy := make([]any, 0, len(lhs)*n)
			for i := 0; i < n; i++ {
				copy = append(copy, lhs...)
			}
			return copy, nil
		}
	}

	return nil, fmt.Errorf("cannot multiply %T and %T", lhs, rhs)
}

func div(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs / rhs, nil
		}
	}

	return nil, fmt.Errorf("cannot divide %T and %T", lhs, rhs)
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

const Ext = ".gqlt"

// Discover all `gqlt` tests in the given directory (recursively).
// Returns the paths of all the test files.
func Discover(dir string) ([]string, error) {
	return doublestar.FilepathGlob(fmt.Sprintf("%s/**/*%s", dir, Ext))
}

type RunOpt func(*runConfig)

// Apply a glob filter to the test files
// This is applied to the path formed by stripping the root from the test file's path.
func WithGlob(glob string) RunOpt {
	return func(o *runConfig) {
		o.glob = glob
	}
}

type runConfig struct {
	// Filter is a glob pattern (with support for **) that is matched against each test file path.
	glob string
}

// `Run` all `gqlt` tests in the given directory (recursively).
func (e *Executor) Run(t *testing.T, ctx context.Context, root string, opts ...RunOpt) {
	runConfig := runConfig{
		glob: "**",
	}

	for _, opt := range opts {
		opt(&runConfig)
	}

	paths, err := Discover(root)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		path := path

		matches, err := doublestar.PathMatch(runConfig.glob, path)
		if err != nil {
			t.Fatal(err)
		}

		if !matches {
			t.SkipNow()
		}

		idx := strings.Index(path, root)
		name := path[idx+len(root)+1 : len(path)-len(Ext)]

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parser, err := parser.NewFromPath(path)
			if err != nil {
				t.Fatal(err)
			}

			file, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			if err := e.runFile(ctx, file); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func (e *Executor) runFile(ctx context.Context, file syn.File) error {
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

		case *syn.SetStmt:
			val, err := e.eval(ctx, ecx, stmt.Value)
			if err != nil {
				return err
			}

			if err := ecx.settings.Set(stmt.Key, val); err != nil {
				return err
			}

		case *syn.AssertStmt:
			bin, ok := stmt.Expr.(*syn.BinaryExpr)
			if ok && bin.Op == lex.Equals2 {
				// special case for common equality assertions to have a better error message
				lhs, err := e.eval(ctx, ecx, bin.Left)
				if err != nil {
					return err
				}

				rhs, err := e.eval(ctx, ecx, bin.Right)
				if err != nil {
					return err
				}

				if !eq(lhs, rhs) {
					return fmt.Errorf("assertion failed: %v != %v", lhs, rhs)
				}
			} else {
				val, err := e.eval(ctx, ecx, stmt.Expr)
				if err != nil {
					return err
				}

				if !truthy(val) {
					return fmt.Errorf("assertion failed")
				}

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
	case *syn.WildcardPat:
		return nil

	case *syn.NamePat:
		binder.bind(pat.Name, val)
		return nil

	case *syn.ListPat:
		vals, ok := val.([]any)
		if !ok {
			return fmt.Errorf("cannot bind %T value `%v` to list pattern", val, val)
		}
		return bindListPat(binder, pat, vals)

	case *syn.ObjectPat:
		vals, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("cannot bind %T value `%v` to object pattern", val, val)
		}
		return bindObjectPat(binder, pat, vals)

	case *syn.LiteralPat:
		if pat.Value != val {
			return fmt.Errorf("literal pattern does not match value: %v != %v", pat.Value, val)
		}
		return nil

	case *syn.RestPat:
		panic("rest pattern should have special handling in list and object cases")

	default:
		panic(fmt.Sprintf("missing pattern bind case: %T", pat))
	}
}

func bindListPat(binder binder, pat *syn.ListPat, values []any) error {
	for i, pat := range pat.Pats {
		rest, ok := pat.(*syn.RestPat)
		if ok {
			if err := bindPat(binder, rest.Pat, values[i:]); err != nil {
				return err
			}
			return nil

		}

		if i > len(values)-1 {
			return bindPat(binder, pat, nil)
		}

		if err := bindPat(binder, pat, values[i]); err != nil {
			return err
		}
	}

	return nil
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

func formatOperation(operation *ast.OperationDefinition) string {
	buf := bytes.NewBufferString("")
	f := formatter.NewFormatter(buf, formatter.WithIndent("  "))
	f.FormatQueryDocument(&ast.QueryDocument{
		Operations: []*ast.OperationDefinition{operation},
	})
	return buf.String()
}

func mapSlice[T, U any](xs []T, f func(T) U) []U {
	ys := make([]U, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func applyNamespaces(namespace []string, operation *ast.OperationDefinition) *ast.OperationDefinition {
	selectionSet := operation.SelectionSet
	// iterate in reverse order to build up the selection set from the inside out
	for i := len(namespace) - 1; i >= 0; i-- {
		selectionSet = ast.SelectionSet{
			&ast.Field{Alias: namespace[i], Name: namespace[i], SelectionSet: selectionSet},
		}
	}

	return &ast.OperationDefinition{
		Operation:           operation.Operation,
		Name:                operation.Name,
		VariableDefinitions: operation.VariableDefinitions,
		Directives:          operation.Directives,
		SelectionSet:        selectionSet,
		Position:            operation.Position,
		Comment:             operation.Comment,
	}
}

func (e *Executor) eval(ctx context.Context, ecx *executionContext, expr syn.Expr) (any, error) {
	switch expr := expr.(type) {
	case *syn.OperationExpr:
		query := formatOperation(applyNamespaces(ecx.settings.namespace, expr.Operation))

		// Pass our local variables directly also as graphql variables
		var data any
		req := Request{Query: query, Variables: ecx.scope.gqlVars()}
		if err := e.client.Request(ctx, req, &data); err != nil {
			return nil, err
		}

		return flatten(data), nil

	case *syn.LiteralExpr:
		return expr.Value, nil

	case *syn.UnaryExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		switch expr.Op {
		case lex.Minus:
			switch lhs := val.(type) {
			case float64:
				return -lhs, nil
			default:
				return nil, fmt.Errorf("cannot negate %T", lhs)
			}
		case lex.Bang, lex.Not:
			return !truthy(val), nil
		default:
			panic(fmt.Sprintf("missing unary expr eval case: %s", expr.Op))
		}

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
			return eq(lhs, rhs), nil
		case lex.BangEqual:
			return !eq(lhs, rhs), nil
		case lex.Plus:
			return add(lhs, rhs)
		case lex.Minus:
			return sub(lhs, rhs)
		case lex.Star:
			return mul(lhs, rhs)
		case lex.Slash:
			return div(lhs, rhs)
		default:
			panic(fmt.Sprintf("missing binary expr eval case: %s", expr.Op))
		}

	case *syn.NameExpr:
		val, ok := ecx.scope.Lookup(expr.Name)
		if !ok {
			return nil, fmt.Errorf("reference to undefined variable: %s", expr.Name)
		}

		return val, nil

	case *syn.ListExpr:
		vals := make([]any, len(expr.Exprs))
		for i, expr := range expr.Exprs {
			val, err := e.eval(ctx, ecx, expr)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}

		return vals, nil

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

	case *syn.IndexExpr:
		val, err := e.eval(ctx, ecx, expr.Expr)
		if err != nil {
			return nil, err
		}

		idx, err := e.eval(ctx, ecx, expr.Index)
		if err != nil {
			return nil, err
		}

		switch val := val.(type) {
		case []any:
			switch idx := idx.(type) {
			case float64:
				if idx < 0 || idx >= float64(len(val)) {
					return nil, nil
				}

				return val[int(idx)], nil
			}
		case map[string]any:
			switch idx := idx.(type) {
			case string:
				if val, ok := val[idx]; ok {
					return val, nil
				}
				return nil, nil
			}
		}

		return nil, fmt.Errorf("cannot index %T with %T", val, idx)

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
