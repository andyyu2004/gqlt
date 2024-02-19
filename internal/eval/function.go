package eval

import (
	"fmt"
	"log"
	"regexp"

	"github.com/google/go-cmp/cmp"
	"github.com/movio/gqlt/gqlparser/ast"
)

func regexMatch(pos ast.HasPosition, lhs, rhs string) (bool, error) {
	regex, err := regexp.Compile(rhs)
	if err != nil {
		return false, errorf(pos, "invalid regex: %v", err)
	}

	return regex.MatchString(lhs), nil
}

func eq(lhs, rhs any) bool {
	return cmp.Equal(lhs, rhs)
}

func add(pos ast.HasPosition, lhs, rhs any) (any, error) {
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

	return nil, errorf(pos, "cannot add %T and %T", lhs, rhs)
}

func sub(pos ast.HasPosition, lhs, rhs any) (any, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs - rhs, nil
		}
	}

	return nil, errorf(pos, "cannot subtract %T and %T", lhs, rhs)
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

// FIXME typecheck these?
var builtinScope = &scope{
	vars: map[string]any{
		"dbg": function(func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("dbg takes exactly 1 argument")
			}
			log.Println(args...)
			return args[0], nil
		}),
		"print": function(func(args []any) (any, error) {
			log.Println(args...)
			return nil, nil
		}),
		"len": function(func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("len takes exactly 1 argument")
			}

			switch arg := args[0].(type) {
			case nil:
				return float64(0), nil
			case string:
				return float64(len(arg)), nil
			case []any:
				return float64(len(arg)), nil
			case map[string]any:
				return float64(len(arg)), nil
			default:
				return nil, fmt.Errorf("len takes a string, array, or object")
			}
		}),
	},
}
