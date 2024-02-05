package gqlt_test

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/andyyu2004/gqlt"

	_ "embed"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//go:embed tests/schema.graphql
var schema string

const testpath = "tests"

func TestGqlt(t *testing.T) {
	q := &query{
		dogs: []dog{
			{
				id:   "1",
				name: "Buddy",
			},
		},
		cats: []cat{
			{
				id:   "1",
				name: "Chips",
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
		factory := gqlt.ClientFactoryFunc(func(testing.TB) (context.Context, gqlt.Client) {
			// reset the counter for each test
			q.counter.Store(0)
			return context.Background(), client
		})
		gqlt.New().Test(t, testpath, factory, gqlt.TypeCheck(true))
	}
}

// convenience function for debugging
// Put whatever in `scratch.gqlt` and debug this test
func TestScratch(t *testing.T) {
	q := &query{}
	schema := graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())
	client := gqlt.GraphQLGophersClient{Schema: schema}
	factory := gqlt.ClientFactoryFunc(func(testing.TB) (context.Context, gqlt.Client) {
		return context.Background(), client
	})
	gqlt.New().Test(t, testpath, factory, gqlt.TypeCheck(true), gqlt.WithGlob("scratch.gqlt"))
}

type AnimalFilter struct {
	Kind *AnimalKind
	Name *string
}

type query struct {
	counter atomic.Int32
	cats    []cat
	dogs    []dog
}

func (q *query) Int() int32 {
	return 1
}

type Foo struct {
	ID      graphql.ID
	String  string
	Int     int32
	Float   float64
	Boolean bool
}

func (q *query) Foo() Foo {
	return Foo{
		ID:      "1",
		String:  "foo",
		Int:     1,
		Float:   1.1,
		Boolean: true,
	}
}

func (q *query) Foos() []Foo {
	return []Foo{
		{
			ID:      "2",
			String:  "3",
			Int:     2,
			Float:   3.3,
			Boolean: false,
		},
	}
}

type Recursive struct {
	ID   graphql.ID
	Next *Recursive
}

func (q *query) Recursive() Recursive {
	return Recursive{
		ID: "1",
		Next: &Recursive{
			ID: "2",
			Next: &Recursive{
				ID: "3",
			},
		},
	}
}

func (q *query) Inc() int32 {
	return q.counter.Add(1)
}

func (q *query) Fail(args struct{ Yes bool }) (int32, error) {
	if args.Yes {
		return 1, errors.New("failed")
	}

	return 0, nil
}
func (q *query) Animals() *query        { return q }
func (q *query) Dogs() dogQuery         { return dogQuery{q} }
func (q *query) Cats() catQuery         { return catQuery{q} }
func (q *query) AllKinds() []AnimalKind { return []AnimalKind{dog{}.Kind(), cat{}.Kind()} }
func (q *query) KindToString(args struct{ Kind AnimalKind }) string {
	return strings.ToLower(string(args.Kind))
}

type Animal struct {
	animal
}

func (a Animal) ToDog() (*dog, bool) { d, ok := a.animal.(*dog); return d, ok }

func (a Animal) ToCat() (*cat, bool) { c, ok := a.animal.(*cat); return c, ok }

type animal interface {
	ID() graphql.ID
	Kind() AnimalKind
	Name() string
}

func (q *query) Search(args struct{ Filter *AnimalFilter }) []Animal {
	var animals []Animal

	for _, dog := range q.dogs {
		if args.Filter == nil {
			continue
		}
		if args.Filter.Kind == nil || *args.Filter.Kind == dog.Kind() {
			if args.Filter.Name == nil || *args.Filter.Name == dog.Name() {
				animals = append(animals, Animal{dog})
			}
		}
	}

	for _, cat := range q.cats {
		if args.Filter == nil {
			continue
		}

		if args.Filter.Kind == nil || *args.Filter.Kind == cat.Kind() {
			if args.Filter.Name == nil || *args.Filter.Name == cat.Name() {
				animals = append(animals, Animal{cat})
			}
		}
	}

	return animals
}

type dogQuery struct{ *query }

func (q dogQuery) First() *dog {
	if len(q.dogs) > 0 {
		return &q.dogs[0]
	}
	return nil
}

func (q dogQuery) List() []dog { return q.dogs }

func (q dogQuery) Find(args struct{ Name string }) *dog {
	for _, dog := range q.dogs {
		if dog.Name() == args.Name {
			return &dog
		}
	}
	return nil
}

type dog struct {
	id   graphql.ID
	name string
}

func (d dog) ID() graphql.ID   { return d.id }
func (d dog) Kind() AnimalKind { return "DOG" }
func (d dog) Name() string     { return d.name }

type AnimalKind string

type cat struct {
	id   graphql.ID
	name string
}

type catQuery struct{ *query }

func (q catQuery) First() *cat {
	if len(q.cats) > 0 {
		return &q.cats[0]
	}
	return nil
}

func (q catQuery) List() []cat { return q.cats }

func (q catQuery) Find(args struct{ Name string }) *cat {
	for _, cat := range q.cats {
		if cat.Name() == args.Name {
			return &cat
		}
	}
	return nil
}

func (c cat) ID() graphql.ID   { return c.id }
func (c cat) Kind() AnimalKind { return "CAT" }
func (c cat) Name() string     { return c.name }
