package lib

import (
	"reflect"
)

func Assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func Ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func Ref[T any](t T) *T {
	return &t
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func TypeOf[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func IsNil(i any) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	return val.Kind() == reflect.Ptr && val.IsNil()
}
