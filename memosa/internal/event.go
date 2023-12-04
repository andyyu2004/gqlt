package internal

import (
	"reflect"
)

type Event interface {
	isEvent()
}

type DidValidateMemoizedValue struct {
	QueryType reflect.Type
	Key       any
}

var _ Event = DidValidateMemoizedValue{}

func (DidValidateMemoizedValue) isEvent() {}

type WillExecute struct {
	QueryType reflect.Type
	Key       any
}

func (WillExecute) isEvent() {}
