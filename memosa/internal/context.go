package internal

import "github.com/andyyu2004/memosa/lib"

type InputKey struct{}

type Input[T any] interface {
	Query[InputKey, T]
}

// might be more ergonomic to have the key implement the query interface and drop the separate K
type Query[K, V any] interface {
	Execute(*Context, K) V
}

type Context struct {
	rt *runtime
}

type Option func(*options)

type options struct {
	eventHandler func(Event)
}

func WithEventHandler(handler func(Event)) Option {
	return func(opts *options) {
		opts.eventHandler = handler
	}
}

func NewContext(opts ...Option) *Context {
	options := options{eventHandler: nil}
	for _, opt := range opts {
		opt(&options)
	}

	return &Context{
		rt: newRt(options.eventHandler),
	}
}

func Set[I Input[T], T any](ctx *Context, value T) {
	set[I, T](ctx.rt, value)
}

func Fetch[Q Query[K, V], K, V any](ctx *Context, key K) V {
	return fetch(ctx, lib.TypeOf[Q](), key).(V)
}
