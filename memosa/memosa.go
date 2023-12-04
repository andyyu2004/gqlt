package memosa

import (
	"github.com/andyyu2004/memosa/internal"
)

type (
	Context                  = internal.Context
	Event                    = internal.Event
	DidValidateMemoizedValue = internal.DidValidateMemoizedValue
	WillExecute              = internal.WillExecute

	Query[K comparable, V comparable] internal.Query[K, V]
	Input[T comparable]               internal.Input[T]
	InputKey                          = internal.InputKey
)

var (
	NewContext       = internal.NewContext
	WithEventHandler = internal.WithEventHandler
)

func Fetch[Q Query[K, V], K comparable, V comparable](ctx *Context, key K) V {
	return internal.Fetch[Q, K, V](ctx, key)
}

func Set[I Input[T], T comparable](ctx *Context, value T) {
	internal.Set[I, T](ctx, value)
}
