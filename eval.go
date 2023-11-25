package gqlt

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/andyyu2004/gqlt/parser"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/bmatcuk/doublestar/v4"
)

type Opt func(*runConfig)

// Apply a glob filter to the test files
// This is applied to the path formed by stripping the root from the test file's path.
func WithGlob(glob string) Opt {
	return func(o *runConfig) {
		o.glob = glob
	}
}

type runConfig struct {
	// Filter is a glob pattern (with support for **) that is matched against each test file path.
	glob string
}

// A thread-safe graphql client
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

func isGqlVar(val any) bool {
	switch reflect.ValueOf(val).Kind() {
	case reflect.Func:
		return false
	default:
		return true
	}
}

func (s *scope) gqlVars() map[string]any {
	vars := map[string]any{}

	for name, val := range s.vars {
		// We avoid passing in function values as graphql variables.
		// This is only a shallow check so we can still pass in a map containing functions for example.
		if !isGqlVar(val) {
			continue
		}

		vars[name] = val
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

const Ext = ".gqlt"

// Discover all `gqlt` tests in the given directory (recursively).
// Returns the paths of all the test files.
func Discover(dir string) ([]string, error) {
	return doublestar.FilepathGlob(fmt.Sprintf("%s/**/*%s", dir, Ext))
}

// `Run` all `gqlt` tests in the given directory (recursively).
func (e *Executor) Run(t *testing.T, ctx context.Context, root string, opts ...Opt) {
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
		if err := e.stmt(ctx, ecx, stmt); err != nil {
			return err
		}
	}

	return nil
}
