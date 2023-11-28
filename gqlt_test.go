package gqlt_test

import (
	"context"
	"strings"
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
		cats: []cat{
			{
				ID:   "1",
				Name: "Chips",
			},
		},
	}

	schema := graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())
	handler := &relay.Handler{Schema: schema}

	clients := []gqlt.Client{
		gqlt.GraphQLGophersClient{Schema: schema},
		gqlt.HTTPClient{Handler: handler},
	}

	for _, client := range clients {
		gqlt.New(client).Run(t, ctx, testpath, gqlt.WithGlob(debugGlob))
	}
}

type query struct {
	cats []cat
	dogs []dog
}

func castSlice[T, U any](xs []T) []U {
	ys := make([]U, len(xs))
	for i, x := range xs {
		ys[i] = any(x).(U)
	}
	return ys
}

func (q query) Animals() query         { return q }
func (q query) Dogs() dogQuery         { return dogQuery{q} }
func (q query) Cats() catQuery         { return catQuery{q} }
func (q query) AllKinds() []AnimalKind { return []AnimalKind{dog{}.Kind(), cat{}.Kind()} }
func (q query) KindToString(args struct{ Kind AnimalKind }) string {
	return strings.ToLower(string(args.Kind))
}

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

type AnimalKind string

func (d dog) Kind() AnimalKind { return "DOG" }

type cat struct {
	ID   graphql.ID
	Name string
}

type catQuery struct{ query }

func (q catQuery) First() *cat {
	if len(q.cats) > 0 {
		return &q.cats[0]
	}
	return nil
}

func (q catQuery) List() []cat { return q.cats }

func (q catQuery) Find(args struct{ Name string }) *cat {
	for _, cat := range q.cats {
		if cat.Name == args.Name {
			return &cat
		}
	}
	return nil
}

func (c cat) Kind() AnimalKind { return "CAT" }
