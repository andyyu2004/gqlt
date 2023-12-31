package lexer_test

import (
	"testing"

	"github.com/andyyu2004/gqlt/gqlparser/gqlerror"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/parser/testrunner"
)

func TestLexer(t *testing.T) {
	testrunner.Test(t, "lexer_test.yml", func(_ *testing.T, input string) testrunner.Spec {
		l := lexer.New(&ast.Source{Input: input, Name: "spec"})

		ret := testrunner.Spec{}
		for {
			tok, err := l.ReadToken()
			if err != nil {
				ret.Error = err.(*gqlerror.Error)
				break
			}

			if tok.Kind == lexer.EOF {
				break
			}

			ret.Tokens = append(ret.Tokens, testrunner.Token{
				Kind:   tok.Kind.Name(),
				Value:  tok.Value,
				Line:   tok.Pos.Line,
				Column: tok.Pos.Column,
				Start:  tok.Pos.Start,
				End:    tok.Pos.End,
				Src:    tok.Pos.Src.Name,
			})
		}

		return ret
	})
}
