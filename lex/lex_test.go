package lex_test

import (
	"testing"

	"github.com/andyyu2004/gqlt/lex"

	"github.com/andyyu2004/gqlt/gqlparser/lexer"
	"github.com/stretchr/testify/require"
)

func TestTokenNumbers(t *testing.T) {
	// Our numbering must match the original tokens as we cast between them assuming this is the case

	tokenKinds := []string{}
	for kind := lex.Invalid; kind <= lex.Comment; kind++ {
		tokenKinds = append(tokenKinds, kind.String())
	}

	tokenTypes := []string{}
	for kind := lexer.Invalid; kind <= lexer.Comment; kind++ {
		tokenTypes = append(tokenTypes, kind.String())
	}

	require.Equal(t, tokenKinds, tokenTypes)
}
