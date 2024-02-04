package eval

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/andyyu2004/gqlt/internal/parser"
	"github.com/andyyu2004/gqlt/internal/typecheck"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/bmatcuk/doublestar/v4"
)

type Opt func(*runConfig)

// Apply a glob filter to the test files
// This is applied to the path formed by stripping the root from the test file's path.
func WithGlob(glob string) Opt {
	return func(c *runConfig) {
		c.glob = glob
	}
}

func TypeCheck(b bool) Opt {
	return func(c *runConfig) {
		c.typecheck = b
	}
}

type runConfig struct {
	// Filter is a glob pattern (with support for **) that is matched against each test file path.
	glob      string
	typecheck bool
}

// A thread-safe graphql client
type Executor struct {
	schemaOnce sync.Once
	schema     schema
}

func New() *Executor {
	return &Executor{}
}

type settings struct {
	namespace []string
}

// keep in sync with typecheck
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
	// path to the file being executed
	path     string
	client   Client
	scope    *scope
	settings settings
}

type scope struct {
	parent *scope
	vars   map[string]any
	// defined fragments that can be referenced by queries
	// map from fragment name to raw fragment string
	fragments map[string]string
}

func (s *scope) bind(name string, val any) {
	s.vars[name] = val
}

func (s *scope) LookupFragment(name string) (string, bool) {
	frag, ok := s.fragments[name]
	if !ok && s.parent != nil {
		return s.parent.LookupFragment(name)
	}

	return frag, ok
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

const Ext = ".gqlt"

// Discover all `gqlt` tests in the given directory (recursively).
// Returns the paths of all the test files.
func Discover(dir string) ([]string, error) {
	return doublestar.FilepathGlob(fmt.Sprintf("%s/**/*%s", dir, Ext))
}

// `Test` all `gqlt` tests in the given directory (recursively).
func (e *Executor) Test(t *testing.T, root string, factory ClientFactory, opts ...Opt) {
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

		idx := strings.Index(path, root)
		name := path[idx+len(root)+1 : len(path)-len(Ext)]

		t.Run(name, func(t *testing.T) {
			if !matches {
				t.SkipNow()
			}

			ctx, client := factory.CreateClient(t)

			if err := e.RunFile(ctx, client, path, opts...); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func (e *Executor) RunFile(ctx context.Context, client Client, uri string, opts ...Opt) error {
	// FIXME, doesn't make much sense to take the `glob` config here as it's not needed
	runConfig := runConfig{}
	for _, opt := range opts {
		opt(&runConfig)
	}

	parser, err := parser.NewFromPath(uri)
	if err != nil {
		return err
	}

	file, err := parser.Parse()
	if err != nil {
		return err
	}

	// No need to pass as schema here for detailed typechecking (the typechecker is aimed at providing errors during editing)
	// We can just fail at runtime.
	tcx := typecheck.New(nil)
	info := tcx.Check(file)
	if len(info.Errors) > 0 && runConfig.typecheck {
		return info.Errors
	}

	if err := e.prepareSchema(ctx, client); err != nil {
		return err
	}

	ecx := &executionContext{
		path:   uri,
		client: client,
		scope: &scope{
			parent:    builtinScope,
			vars:      map[string]any{},
			fragments: map[string]string{},
		},
	}

	for _, stmt := range file.Stmts {
		if err := e.stmt(ctx, ecx, stmt); err != nil {
			return err
		}
	}

	return nil
}

const introspectionQuery = `
query {
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    types {
      kind
      name
      inputFields {
        name
        type {
          ...TypeRef
        }
      }
      fields(includeDeprecated: true) {
        name
        args {
          name
          type {
            ...TypeRef
          }
        }
        type {
          ...TypeRef
        }
      }
    }
  }
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
`

type schema struct {
	QueryType    typename
	MutationType typename
	Types        map[typename]ty
}

type typename string

type ty struct {
	Name        typename
	Kind        syn.DefinitionKind
	InputFields map[string]field
	Fields      map[string]field
}

type tyref struct {
	Kind   syn.DefinitionKind
	Name   typename
	OfType *tyref
}

func (t tyref) LeafType() typename {
	if t.OfType == nil {
		// leaf-types should have a non-empty name
		lib.Assert(t.Name != "")
		return t.Name
	}

	return t.OfType.LeafType()
}

type field struct {
	Type typename
	Args map[string]typename
}

func (e *Executor) prepareSchema(ctx context.Context, client Client) error {
	var err error
	e.schemaOnce.Do(func() {
		type Field struct {
			Name string
			Args []struct {
				Name string
				Type tyref
			}
			Type tyref
		}
		var res struct {
			Schema struct {
				QueryType struct {
					Name typename
				}
				MutationType struct {
					Name typename
				}
				Types []struct {
					Kind        syn.DefinitionKind
					Name        typename
					Fields      []Field
					InputFields []Field
				} `json:"types"`
			} `json:"__schema"`
		}

		var errors GraphQLErrors
		errors, err = client.Request(ctx, Request{Query: introspectionQuery}, &res)
		if errors != nil {
			err = errors
			return
		}

		// can continue even on error safely, it will just become mostly a noop
		types := map[typename]ty{}
		for _, t := range res.Schema.Types {
			transformFields := func(fields []Field) map[string]field {
				out := make(map[string]field, len(fields))
				for _, f := range fields {
					args := make(map[string]typename, len(f.Args))
					for _, arg := range f.Args {
						args[arg.Name] = arg.Type.LeafType()
					}
					out[f.Name] = field{Args: args, Type: f.Type.LeafType()}
				}
				return out
			}

			fields := transformFields(t.Fields)
			inputFields := transformFields(t.InputFields)
			types[t.Name] = ty{Name: t.Name, Kind: t.Kind, Fields: fields, InputFields: inputFields}
		}

		e.schema = schema{
			QueryType:    res.Schema.QueryType.Name,
			MutationType: res.Schema.MutationType.Name,
			Types:        types,
		}
	})
	return err
}

// an interface an error type should implement if it should be catchable by a `try <expr>`
type catchable interface {
	// convert the error into a gqlt value that can returned
	catch() any
}
