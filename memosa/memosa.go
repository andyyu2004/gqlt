package memosa

import (
	"github.com/movio/gqlt/memosa/internal"
)

type (
	Context                  = internal.Context
	Event                    = internal.Event
	DidValidateMemoizedValue = internal.DidValidateMemoizedValue
	WillExecute              = internal.WillExecute

	Query[K, V any] internal.Query[K, V]
	Input[T any]    internal.Input[T]
	InputKey        = internal.InputKey
)

var (
	New              = internal.NewContext
	WithEventHandler = internal.WithEventHandler
)

func Get[I Input[T], T any](ctx *Context) T {
	return Fetch[I, InputKey, T](ctx, InputKey{})
}

func Fetch[Q Query[K, V], K, V any](ctx *Context, key K) V {
	return internal.Fetch[Q, K, V](ctx, key)
}

func Set[I Input[T], T any](ctx *Context, value T) {
	internal.Set[I, T](ctx, value)
}
