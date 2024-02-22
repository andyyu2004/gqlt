package eval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

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

func lt(lhs, rhs any) (bool, error) {
	switch lhs := lhs.(type) {
	case float64:
		switch rhs := rhs.(type) {
		case float64:
			return lhs < rhs, nil
		}
	case string:
		switch rhs := rhs.(type) {
		case string:
			return lhs < rhs, nil
		}
	}

	return false, fmt.Errorf("cannot compare %T and %T", lhs, rhs)
}

func lte(lhs, rhs any) (bool, error) {
	if eq(lhs, rhs) {
		return true, nil
	}

	return lt(lhs, rhs)
}

func gt(lhs, rhs any) (bool, error) {
	lte, err := lte(lhs, rhs)
	return !lte, err
}

func gte(lhs, rhs any) (bool, error) {
	lt, err := lt(lhs, rhs)
	return !lt, err
}

func fetch(url string, args map[string]any) (any, error) {
	method := "GET"
	if m, ok := args["method"]; ok {
		if s, ok := m.(string); ok {
			method = s
		} else {
			return nil, fmt.Errorf("method must be a string")
		}
	}

	var body io.Reader
	if b, ok := args["body"]; ok {
		byts, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(byts)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if headers, ok := args["headers"]; ok {
		h, ok := headers.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("headers must be an object")
		}

		for k, v := range h {
			switch v := v.(type) {
			case string:
				req.Header.Set(k, v)
			default:
				return nil, fmt.Errorf("header values must be strings")
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	headers := map[string][]any{}
	for k, v := range resp.Header {
		values := make([]any, len(v))
		for i, s := range v {
			values[i] = s
		}
		headers[k] = values
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"statusCode": float64(resp.StatusCode),
		"status":     resp.Status,
		"headers":    headers,
		"body":       string(data),
	}, nil
}

// FIXME typecheck these?
var builtinScope = &scope{
	vars: map[string]any{
		"now": function(func(args []any) (any, error) {
			if len(args) > 1 {
				return nil, fmt.Errorf("now() or now(format: string) expected")
			}

			format := time.RFC3339
			if len(args) == 1 {
				switch arg := args[0].(type) {
				case string:
					format = arg
				default:
					return nil, fmt.Errorf("format argument must be a string")
				}
			}
			return time.Now().Format(format), nil
		}),
		"fetch": function(func(args []any) (any, error) {
			const help = `fetch(url: string, args: {
				method?: string,
				body:? string,
				headers: { string: string },
			})`
			if len(args) < 1 || len(args) > 2 {
				return nil, fmt.Errorf(help)
			}

			url, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("%s: 'url' must be a string", help)
			}

			if len(args) == 1 {
				return fetch(url, map[string]any{})
			} else {
				p, ok := args[1].(map[string]any)
				if !ok {
					return nil, fmt.Errorf("%s: 'args' must be an object", help)
				}
				return fetch(url, p)
			}
		}),
		"parseJSON": function(func(args []any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("parseJson takes exactly 1 argument")
			}

			switch arg := args[0].(type) {
			case string:
				var val any
				if err := json.Unmarshal([]byte(arg), &val); err != nil {
					return nil, err
				}
				return val, nil
			default:
				return nil, fmt.Errorf("parseJSON takes a string")
			}
		}),
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
