package internal

import (
	"fmt"
	"reflect"
)

func eq(lhs, rhs any) bool {
	return reflect.DeepEqual(lhs, rhs)
}

func add(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs + rhs, nil
		}
	case string:
		switch rhs := rhs.(type) {
		case string:
			return lhs + rhs, nil
		}
	case []any:
		switch rhs := rhs.(type) {
		case []any:
			return append(lhs, rhs...), nil
		}
	}

	return nil, fmt.Errorf("cannot add %T and %T", lhs, rhs)
}

func sub(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs - rhs, nil
		}
	}

	return nil, fmt.Errorf("cannot subtract %T and %T", lhs, rhs)
}

func mul(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs * rhs, nil
		}
	case []any:
		switch rhs := rhs.(type) {
		case float64:
			n := int(rhs)
			copy := make([]any, 0, len(lhs)*n)
			for i := 0; i < n; i++ {
				copy = append(copy, lhs...)
			}
			return copy, nil
		}
	}

	return nil, fmt.Errorf("cannot multiply %T and %T", lhs, rhs)
}

func div(lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs / rhs, nil
		}
	}

	return nil, fmt.Errorf("cannot divide %T and %T", lhs, rhs)
}

func truthy(val any) bool {
	switch val := val.(type) {
	case nil:
		return false
	case bool:
		return val
	case int:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

var builtinScope = &scope{
	vars: map[string]any{
		"example": function(func(args []any) (any, error) {
			if err := checkArity(1, args); err != nil {
				return nil, err
			}

			if !truthy(args[0]) {
				return nil, fmt.Errorf("assertion failed")
			}

			return nil, nil
		}),
	},
}
