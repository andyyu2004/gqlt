package gqlt_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andyyu2004/gqlt"
	"github.com/andyyu2004/gqlt/parser"

	_ "embed"

	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/require"
)

//go:embed tests/schema.graphql
var schema string

// builtin glob doesn't implement ** :/
func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func TestGqlt(t *testing.T) {
	// set this to debug a particular case i.e. `binary.gqlt`
	const debugFilter = "binary.gqlt"

	ctx := context.Background()

	const testpath = "tests"

	paths, err := glob(testpath, ".gqlt")
	require.NoError(t, err)

	q := &query{
		dogs: []dog{
			{
				ID:   "1",
				Name: "Buddy",
			},
		},
	}
	client := schemaClientAdaptor{graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())}

	for _, path := range paths {
		path := path

		if !strings.HasSuffix(path, debugFilter) {
			t.SkipNow()
		}

		idx := strings.Index(path, testpath)
		require.True(t, idx != -1)
		name := path[idx+len(testpath)+1:]

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parser, err := parser.NewFromPath(path)
			require.NoError(t, err)

			file, err := parser.Parse()
			require.NoError(t, err)

			executor := gqlt.New(client)
			require.NoError(t, executor.Run(ctx, file))
		})
	}
}

type schemaClientAdaptor struct {
	schema *graphql.Schema
}

func (a schemaClientAdaptor) Request(ctx context.Context, query string, variables map[string]any, out any) error {
	res := a.schema.Exec(ctx, query, "", variables)
	if len(res.Errors) > 0 {
		errs := make([]error, 0, len(res.Errors))
		for _, err := range res.Errors {
			errs = append(errs, err)
		}

		return errors.Join(errs...)
	}

	return json.Unmarshal([]byte(res.Data), out)
}

var _ gqlt.Client = schemaClientAdaptor{}

type query struct {
	dogs []dog
}

func (q query) Animals() query { return q }
func (q query) Dogs() dogQuery { return dogQuery{q} }

type dogQuery struct{ query }

func (q dogQuery) First() *dog {
	if len(q.dogs) > 0 {
		return &q.dogs[0]
	}
	return nil
}

func (q dogQuery) List() []dog { return q.dogs }

func (q dogQuery) Find(args struct{ Name string }) *dog {
	for _, dog := range q.dogs {
		if dog.Name == args.Name {
			return &dog
		}
	}
	return nil
}

type dog struct {
	ID   graphql.ID
	Name string
}
