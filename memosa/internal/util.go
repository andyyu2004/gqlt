package internal

import (
	"reflect"
)

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func ref[T any](t T) *T {
	return &t
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func typeof[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}
