package ide

import (
	"fmt"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/iterator"
	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Highlights []Highlight

func (hs Highlights) String() string {
	sb := strings.Builder{}
	for _, h := range hs {
		// we can assume tokens are not multiline
		len := h.Pos.End - h.Pos.Start
		sb.WriteString(fmt.Sprintf("%d:%d..%d:%d: %v\n", h.Pos.Line, h.Pos.Column, h.Pos.Line, h.Pos.Column+len, h.TokenKind))
	}
	return sb.String()
}

type Highlight struct {
	Pos       ast.Position
	TokenKind protocol.SemanticTokenType
}

func (ide *IDE) Highlight(path string) Highlights {
	root := ide.Parse(path)
	return iterator.FilterMap(traverseTokens(root), func(token syn.Token) (Highlight, bool) {
		var kind protocol.SemanticTokenType
		switch token.Kind {
		case lex.Let, lex.False, lex.True, lex.Null:
			kind = protocol.SemanticTokenTypeKeyword
		case lex.Name:
			kind = protocol.SemanticTokenTypeVariable
		case lex.Equals:
			kind = protocol.SemanticTokenTypeOperator
		case lex.Int, lex.Float:
			kind = protocol.SemanticTokenTypeNumber
		case lex.String, lex.BlockString:
			kind = protocol.SemanticTokenTypeString
		default:
			return Highlight{}, false
		}

		return Highlight{Pos: token.Position, TokenKind: kind}, true
	}).Collect()
}

func traverseTokens(node syn.Node) iterator.Iterator[syn.Token] {
	return iterator.FilterMap(syn.Traverse(node), func(child syn.Child) (syn.Token, bool) {
		token, ok := child.(syn.Token)
		return token, ok
	})
}
