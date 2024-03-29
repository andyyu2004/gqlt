package validator

import (
	"testing"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/movio/gqlt/syn"
)

func Test_sameArguments(t *testing.T) {
	tests := map[string]struct {
		args   func() (args1, args2 []*syn.Argument)
		result bool
	}{
		"both argument lists empty": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return nil, nil
			},
			result: true,
		},
		"args 1 empty, args 2 not": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return nil, []*syn.Argument{
					{
						Name:     lexer.Token{Value: "thing"},
						Value:    &syn.Value{Raw: "a thing"},
						Position: ast.Position{},
					},
				}
			},
			result: false,
		},
		"args 2 empty, args 1 not": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return []*syn.Argument{
					{
						Name:     lexer.Token{Value: "thing"},
						Value:    &syn.Value{Raw: "a thing"},
						Position: ast.Position{},
					},
				}, nil
			},
			result: false,
		},
		"args 1 mismatches args 2 names": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return []*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
					},
					[]*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing2"},
							Value:    &syn.Value{Raw: "2 thing"},
							Position: ast.Position{},
						},
					}
			},
			result: false,
		},
		"args 1 mismatches args 2 values": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return []*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
					},
					[]*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "2 thing"},
							Position: ast.Position{},
						},
					}
			},
			result: false,
		},
		"args 1 matches args 2 names and values": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return []*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
					},
					[]*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
					}
			},
			result: true,
		},
		"args 1 matches args 2 names and values where multiple exist in various orders": {
			args: func() (args1 []*syn.Argument, args2 []*syn.Argument) {
				return []*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
						{
							Name:     lexer.Token{Value: "thing2"},
							Value:    &syn.Value{Raw: "2 thing"},
							Position: ast.Position{},
						},
					},
					[]*syn.Argument{
						{
							Name:     lexer.Token{Value: "thing1"},
							Value:    &syn.Value{Raw: "1 thing"},
							Position: ast.Position{},
						},
						{
							Name:     lexer.Token{Value: "thing2"},
							Value:    &syn.Value{Raw: "2 thing"},
							Position: ast.Position{},
						},
					}
			},
			result: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args1, args2 := tc.args()

			resp := sameArguments(args1, args2)

			if resp != tc.result {
				t.Fatalf("Expected %t got %t", tc.result, resp)
			}
		})
	}
}
