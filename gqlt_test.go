package gqlt_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/movio/gqlt"
	"github.com/movio/gqlt/gqlparser"
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/annotate"
	"github.com/stretchr/testify/require"

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

	gqlparserSchema := gqlparser.MustLoadSchema(&ast.Source{Name: "schema.graphql", Input: schema})
	schema := graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())
	handler := &relay.Handler{Schema: schema}

	clients := []gqlt.Client{
		gqlt.GraphQLGophersClient{Schema: schema},
		gqlt.HTTPClient{Handler: handler},
	}

	var lock sync.Mutex
	alreadyReportedError := map[string]struct{}{}

	for _, client := range clients {
		factory := gqlt.ClientFactoryFunc(func(testing.TB) (context.Context, gqlt.Client) {
			// reset the counter for each test
			q.counter.Store(0)
			return context.Background(), client
		})
		gqlt.New().Test(
			t,
			testpath,
			factory,
			gqlt.TypeCheck(true),
			gqlt.WithSchema(gqlparserSchema),
			gqlt.WithErrorHandler(func(t *testing.T, path string, evalErr error) {
				lock.Lock()
				// since we run once per client
				if _, ok := alreadyReportedError[path]; ok {
					defer lock.Unlock()
					return
				}
				alreadyReportedError[path] = struct{}{}
				lock.Unlock()

				bytes, err := os.ReadFile(path)
				require.NoError(t, err)
				annotation := evalErr.(annotate.Annotation)
				annotated := annotate.Annotate(string(bytes), []annotate.Annotation{annotation})
				snaps.WithConfig(snaps.Filename(strings.TrimSuffix(path, ".gqlt"))).MatchSnapshot(t, annotated)
			}),
		)
	}
}

// convenience function for debugging
// Put whatever in `scratch.gqlt` and debug this test
func TestScratch(t *testing.T) {
	q := &query{}
	gqlparserSchema := gqlparser.MustLoadSchema(&ast.Source{Name: "schema.graphql", Input: schema})
	schema := graphql.MustParseSchema(schema, q, graphql.UseFieldResolvers())
	client := gqlt.GraphQLGophersClient{Schema: schema}
	factory := gqlt.ClientFactoryFunc(func(testing.TB) (context.Context, gqlt.Client) {
		return context.Background(), client
	})
	gqlt.New().Test(t, testpath, factory, gqlt.TypeCheck(true), gqlt.WithGlob("**/scratch.gqlt"), gqlt.WithSchema(gqlparserSchema))
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

type Any string

func (Any) ImplementsGraphQLType(name string) bool {
	return name == "Any"
}

func (Any) UnmarshalGraphQL(input interface{}) error {
	return nil
}

type Foo struct {
	Any     Any
	ID      graphql.ID
	String  string
	Int     int32
	Float   float64
	Boolean bool
}

func (q *query) Foo() Foo {
	return Foo{
		Any:     "any",
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

func (q *query) Unions() unions { return unions{q} }

type unions struct{ *query }

type ab struct {
	any
}

func (ab ab) ToA() (*a, bool) { a, ok := ab.any.(a); return &a, ok }
func (ab ab) ToB() (*b, bool) { b, ok := ab.any.(b); return &b, ok }

type a struct {
	ID graphql.ID
	A  int32
}

type b struct {
	ID graphql.ID
	B  bool
}

func (u unions) AB(args struct{ Pick bool }) ab {
	if args.Pick {
		return ab{b{"2", true}}
	}
	return ab{a{"1", 42}}
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
func (d dog) Bark() string     { return "woof" }

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
func (c cat) Meow() string     { return "meow" }
