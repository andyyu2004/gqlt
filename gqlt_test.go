package gqlt_test

import (
	"context"
	"testing"

	"github.com/andyyu2004/gqlt"

	_ "embed"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//go:embed tests/schema.graphql
var schema string

func TestGqlt(t *testing.T) {
	// change this to something else to debug a particular test
	const debugGlob = "**"

	ctx := context.Background()

	const testpath = "tests"

	q := &query{
		dogs: []dog{
			{
				ID:   "1",
				Name: "Buddy",
			},
		},
	}

	schema := graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())
	handler := &relay.Handler{Schema: schema}

	clients := []gqlt.Client{
		gqlt.GraphQLGophersClient{schema},
		gqlt.HTTPClient{Handler: handler},
	}

	for _, client := range clients {
		gqlt.New(client).Run(t, ctx, testpath, gqlt.WithGlob(debugGlob))
	}
}

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
