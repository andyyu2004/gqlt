package gqlt_test

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"andyyu2004/gqlt"
	"andyyu2004/gqlt/parser"

	_ "embed"

	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/require"
)

//go:embed tests/schema.graphql
var schema string

func TestGqlt(t *testing.T) {
	ctx := context.Background()

	paths, err := filepath.Glob(filepath.Join("tests", "**.gqlt"))
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
		_, filename := filepath.Split(path)
		name := filename[:len(filename)-len(filepath.Ext(path))]

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
