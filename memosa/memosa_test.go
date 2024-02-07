package memosa_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/memosa"

	"github.com/stretchr/testify/require"
)

//      InputA     InputB
//       /            \
//   QueryA(X)      QueryC
//    /                \
// QueryB(X)          QueryD

type InputA struct{}

var _ memosa.Input[int32] = InputA{}

func (InputA) Execute(ctx *memosa.Context, _ memosa.InputKey) int32 { panic(0) }

type InputB struct{}

var _ memosa.Input[int32] = InputB{}

func (InputB) Execute(ctx *memosa.Context, _ memosa.InputKey) int32 { panic(0) }

type QueryA struct{}

var _ memosa.Query[KeyA, int32] = QueryA{}

type KeyA struct{ X int32 }

func (q QueryA) Execute(ctx *memosa.Context, key KeyA) int32 {
	return memosa.Get[InputA](ctx) + key.X
}

type QueryB struct{}

var _ memosa.Query[KeyB, int32] = QueryB{}

type KeyB struct{ X int32 }

func (QueryB) Execute(ctx *memosa.Context, key KeyB) int32 {
	return memosa.Fetch[QueryA](ctx, KeyA(key)) + 42
}

type QueryC struct{}

var _ memosa.Query[KeyC, bool] = QueryC{}

type KeyC struct{}

func (QueryC) Execute(ctx *memosa.Context, key KeyC) bool {
	b := memosa.Get[InputB](ctx)
	if b%2 == 0 {
		return true
	} else {
		return false
	}
}

type QueryD struct{}

var _ memosa.Query[KeyD, int32] = QueryD{}

type KeyD struct{}

func (QueryD) Execute(ctx *memosa.Context, key KeyD) int32 {
	c := memosa.Fetch[QueryC](ctx, KeyC{})
	if c {
		return 42
	} else {
		return -1
	}
}

func NewContext() (*memosa.Context, <-chan memosa.Event) {
	ch := make(chan memosa.Event, 1000)
	return memosa.New(memosa.WithEventHandler(func(event memosa.Event) { ch <- event })), ch
}

func eq[T any](t testing.TB, x, y T) {
	t.Helper()
	require.Equal(t, x, y)
}

func fetch[Q memosa.Query[K, V], K, V comparable](t *testing.T, ctx *memosa.Context, key K, expectedValue V, ch <-chan memosa.Event, expectation expect.Expectation) {
	t.Helper()
	eq(t, expectedValue, memosa.Fetch[Q](ctx, key))
	expectation.AssertEqual(t, formatEvents(ch))
}

func formatEvents(ch <-chan memosa.Event) string {
	builder := strings.Builder{}
	for len(ch) > 0 {
		event := <-ch
		builder.WriteString(fmt.Sprintf("\n%s%v", reflect.TypeOf(event).Name(), event))
	}
	return builder.String()
}

func TestSmoke(t *testing.T) {
	ctx, ch := NewContext()
	memosa.Set[InputA](ctx, 2)
	memosa.Set[InputB](ctx, 12)

	fetch[QueryA](t, ctx, KeyA{X: 1}, 3, ch, expect.Expect(`
WillExecute{memosa_test.QueryA {1}}`))
	// same query should be memoized
	fetch[QueryA](t, ctx, KeyA{X: 1}, 3, ch, expect.Expect(`
DidValidateMemoizedValue{memosa_test.QueryA {1}}`))

	fetch[QueryA](t, ctx, KeyA{X: 2}, 4, ch, expect.Expect(`
WillExecute{memosa_test.QueryA {2}}`))
	fetch[QueryB](t, ctx, KeyB{X: 2}, 46, ch, expect.Expect(`
WillExecute{memosa_test.QueryB {2}}
DidValidateMemoizedValue{memosa_test.QueryA {2}}`))

	fetch[QueryC](t, ctx, KeyC{}, true, ch, expect.Expect(`
WillExecute{memosa_test.QueryC {}}`))

	fetch[QueryD](t, ctx, KeyD{}, 42, ch, expect.Expect(`
WillExecute{memosa_test.QueryD {}}
DidValidateMemoizedValue{memosa_test.QueryC {}}`))

	// updating input should invalidate memoized values that depend on it
	memosa.Set[InputA](ctx, 102)

	fetch[QueryA](t, ctx, KeyA{X: 1}, 103, ch, expect.Expect(`
WillExecute{memosa_test.QueryA {1}}`))

	// however, QueryC shouldn't be invalidated as it does not depend on InputA
	fetch[QueryC](t, ctx, KeyC{}, true, ch, expect.Expect(`
DidValidateMemoizedValue{memosa_test.QueryC {}}`))

	// similarly QueryD should not be invalidated
	fetch[QueryD](t, ctx, KeyD{}, 42, ch, expect.Expect(`
DidValidateMemoizedValue{memosa_test.QueryC {}}
DidValidateMemoizedValue{memosa_test.QueryD {}}`))

	memosa.Set[InputB](ctx, 13)

	// needs reexecution since a relevant input has changed
	fetch[QueryD](t, ctx, KeyD{}, -1, ch, expect.Expect(`
WillExecute{memosa_test.QueryC {}}
WillExecute{memosa_test.QueryD {}}
DidValidateMemoizedValue{memosa_test.QueryC {}}`))

	memosa.Set[InputB](ctx, 15)

	// however, we can short circuit before we execute D again because C's output has not changed
	fetch[QueryD](t, ctx, KeyD{}, -1, ch, expect.Expect(`
WillExecute{memosa_test.QueryC {}}
DidValidateMemoizedValue{memosa_test.QueryD {}}`))

	memosa.Set[InputB](ctx, 16)
	fetch[QueryD](t, ctx, KeyD{}, 42, ch, expect.Expect(`
WillExecute{memosa_test.QueryC {}}
WillExecute{memosa_test.QueryD {}}
DidValidateMemoizedValue{memosa_test.QueryC {}}`))
}

func TestQueryIsReexecutedEvenWhenInputsAreUpToDate(t *testing.T) {
	// Even if a query `Q`'s inputs are all up to date, this is not sufficient
	// to say that `Q` is up to date.
	// It's possible that Q's dependencies have been reexecuted since the last input change, but not Q
	ctx, ch := NewContext()
	memosa.Set[InputA](ctx, 2)

	// populate a memo for QueryB
	fetch[QueryB](t, ctx, KeyB{X: 1}, 45, ch, expect.Expect(`
WillExecute{memosa_test.QueryB {1}}
WillExecute{memosa_test.QueryA {1}}`))

	// update inputA to a different value
	memosa.Set[InputA](ctx, 3)

	// compute A directly
	fetch[QueryA](t, ctx, KeyA{X: 1}, 4, ch, expect.Expect(`
WillExecute{memosa_test.QueryA {1}}`))

	// shuld be memoized
	fetch[QueryA](t, ctx, KeyA{X: 1}, 4, ch, expect.Expect(`
DidValidateMemoizedValue{memosa_test.QueryA {1}}`))

	// QueryB must be reexecuted because QueryA has changed value since the last time QueryB was executed
	fetch[QueryB](t, ctx, KeyB{X: 1}, 46, ch, expect.Expect(`
WillExecute{memosa_test.QueryB {1}}
DidValidateMemoizedValue{memosa_test.QueryA {1}}`))

	// update InputA again
	memosa.Set[InputA](ctx, 4)

	// compute QueryA to keep it up to date
	fetch[QueryA](t, ctx, KeyA{X: 1}, 5, ch, expect.Expect(`
WillExecute{memosa_test.QueryA {1}}`))

	// update some other input to bump the revision
	memosa.Set[InputB](ctx, 12)

	// B still needs to reexecute
	fetch[QueryB](t, ctx, KeyB{X: 1}, 46, ch, expect.Expect(`
DidValidateMemoizedValue{memosa_test.QueryA {1}}
DidValidateMemoizedValue{memosa_test.QueryB {1}}`))
}
