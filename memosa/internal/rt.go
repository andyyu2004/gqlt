package internal

import (
	"log"
	"reflect"

	"github.com/movio/gqlt/memosa/internal/hash"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/memosa/stack"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type rev int

type activeQuery struct {
	dependencies deps
}

type deps struct {
	// map from the query type to the keys that were read
	inputs map[reflect.Type][]any
	// the highest revision of any input that was read during the execution of this query
	maxRev rev
}

func newActiveQuery() *activeQuery {
	return &activeQuery{deps{map[reflect.Type][]any{}, 0}}
}

func recordRead(rt *runtime, queryType reflect.Type, key any) {
	a, ok := rt.activeQueryStack.Peek()
	if ok {
		a.dependencies.inputs[queryType] = append(a.dependencies.inputs[queryType], key)
		if input, ok := rt.inputStorages[queryType]; ok {
			a.dependencies.maxRev = max(a.dependencies.maxRev, input.rev)
		} else if query, ok := rt.queryStorages[queryType]; ok {
			memo, ok := query.memoMap.Get(hash.Hash(key))
			lib.Assert(ok)
			a.dependencies.maxRev = max(a.dependencies.maxRev, memo.deps.maxRev)
		}
	}
}

type runtime struct {
	eventHandler     func(Event)
	revision         rev
	activeQueryStack *stack.Stack[*activeQuery]
	inputStorages    map[reflect.Type]*inputStorage
	queryStorages    map[reflect.Type]*queryStorage
}

type inputStorage struct {
	rev   rev
	value any
}

func (ctx *Context) verifyMemo(memo *memoized[any]) bool {
	lib.Assert(memo.verifiedAt <= ctx.rt.revision)
	if memo.verifiedAt == ctx.rt.revision {
		return true
	}

	for queryType, keys := range memo.deps.inputs {
		for _, key := range keys {
			if !ctx.verifyQuery(queryType, key) {
				return false
			}
		}
	}

	lib.Assert(memo.verifiedAt < ctx.rt.revision)
	memo.verifiedAt = ctx.rt.revision
	return true
}

// verifyQuery that the memoized value for the given key is up-to-date, and reexecutes (recursively) if necessary
// returns true if no downstream queries need to be reexecuted
func (ctx *Context) verifyQuery(queryType reflect.Type, key any) bool {
	if _, ok := key.(InputKey); ok {
		// if the key is an input key, then we just check the revision of the input
		return ctx.rt.inputStorages[queryType].rev < ctx.rt.revision
	}

	storage := ctx.rt.queryStorages[queryType]
	memo, ok := storage.memoMap.Get(hash.Hash(key))
	lib.Assert(ok) // we wouldn't end up here if the key wasn't in the map

	if ctx.verifyMemo(memo) {
		if memo.deps.maxRev < ctx.rt.revision {
			ctx.rt.event(DidValidateMemoizedValue{queryType, key})
			return true
		} else {
			return false
		}
	}

	prevMaxRev := memo.deps.maxRev
	prevValue := memo.value
	newValue := execute(ctx, queryType, memo, key)

	if reflect.DeepEqual(newValue, prevValue) {
		lib.Assert(memo.deps.maxRev >= prevMaxRev)
		// Important to set the value back to the previous value,
		// even if they're deeply equal, we must retain identity too.
		memo.value = prevValue
		memo.deps.maxRev = prevMaxRev
		return true
	}

	return false
}

func (rt *runtime) event(event Event) {
	if rt.eventHandler != nil {
		(rt.eventHandler)(event)
	}
}

// storage for a particular query
type queryStorage struct {
	// a lru cache from hashed key to value
	memoMap *simplelru.LRU[hash.Hashed, *memoized[any]]
}

type memoized[T any] struct {
	// the revision at which this memoized value was last verified
	verifiedAt rev
	// the queries that this memoized value depends on
	deps deps
	// the memoized value
	value T
}

func newRt(eventHandler func(Event)) *runtime {
	return &runtime{
		revision:         0,
		activeQueryStack: new(stack.Stack[*activeQuery]),
		inputStorages:    make(map[reflect.Type]*inputStorage),
		queryStorages:    make(map[reflect.Type]*queryStorage),
		eventHandler:     eventHandler,
	}
}

// Get the value of an input. Panics if the input is not set.
func get(ctx *Context, inputType reflect.Type) any {
	storage, ok := ctx.rt.inputStorages[inputType]
	if !ok {
		log.Panicf("input %v not set", inputType)
	}

	return storage.value
}

// Set the value of an input.
func set[I Input[T], T any](rt *runtime, value T) {
	// this optimization is interfering with some correctness tests for now
	// if storage, ok := rt.inputStorages[lib.TypeOf[I]()]; ok && reflect.DeepEqual(storage.value, value) {
	// 	return
	// }

	rt.revision++
	rt.inputStorages[lib.TypeOf[I]()] = &inputStorage{rt.revision, value}
}

func tryGet(ctx *Context, queryType reflect.Type) (any, bool) {
	method, ok := queryType.MethodByName("Execute")
	lib.Assert(ok)

	if method.Type.In(2) == lib.TypeOf[InputKey]() {
		return get(ctx, queryType), true
	}

	return nil, false
}

func fetch(ctx *Context, queryType reflect.Type, key any) any {
	recordRead(ctx.rt, queryType, key)

	// if the query is an input query, return the input value that was set
	if value, ok := tryGet(ctx, queryType); ok {
		return value
	}

	// otherwise, this is a derived query

	// if a up-to-date memoized value exists, return it
	memo, ok := tryFetchMemoized(ctx.rt, queryType, key)
	if ok {
		ctx.rt.event(DidValidateMemoizedValue{queryType, key})
		return memo.value
	}

	if memo == nil {
		// no memoized value exists, must execute the query
		return execute(ctx, queryType, memo, key)
	}

	// otherwise walk the dependency graph and reexecute as necessary
	ctx.verifyQuery(queryType, key)

	return memo.value
}

func execute(ctx *Context, queryType reflect.Type, memo *memoized[any], key any) any {
	if value, ok := tryGet(ctx, queryType); ok {
		// input queries should not have an associated memo
		lib.Assert(memo == nil)
		return value
	}

	ctx.rt.activeQueryStack.Push(newActiveQuery())
	ctx.rt.event(WillExecute{queryType, key})
	value := reflect.New(queryType).MethodByName("Execute").Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(key)})[0].Interface()
	activeQuery := ctx.rt.activeQueryStack.MustPop()

	if memo == nil {
		memoize(ctx.rt, activeQuery.dependencies, queryType, key, value)
	} else {
		memo.deps = activeQuery.dependencies
		memo.value = value
		memo.verifiedAt = ctx.rt.revision
	}

	return value
}

// returns `nil, false` if no memoized value exists
// returns `nil, true` if a memoized value exists but is potentially stale
// returns `memo, true` if a memoized value exists and is definitely up-to-date
func tryFetchMemoized(rt *runtime, queryType reflect.Type, key any) (*memoized[any], bool) {
	storage, ok := rt.queryStorages[queryType]
	if !ok {
		return nil, false
	}

	memo, ok := storage.memoMap.Get(hash.Hash(key))
	if !ok {
		return nil, false
	}

	if memo.verifiedAt != rt.revision {
		return memo, false
	}

	return memo, true
}

const lruSize = 128

func memoize(rt *runtime, deps deps, queryType reflect.Type, key any, value any) {
	storage, ok := rt.queryStorages[queryType]
	if !ok {
		storage = &queryStorage{memoMap: lib.Must(simplelru.NewLRU[hash.Hashed, *memoized[any]](lruSize, nil))}
		rt.queryStorages[queryType] = storage
	}

	storage.memoMap.Add(hash.Hash(key), &memoized[any]{rt.revision, deps, value})
}
