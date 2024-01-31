package ide

import (
	"fmt"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/iterator"
	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/memosa/stack"
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
	type Scope int
	const (
		ScopeObject Scope = iota
	)
	scopes := stack.Stack[Scope]{}
	return iterator.FilterMap(syn.Traverse(root), func(event syn.Event) (Highlight, bool) {
		switch event := event.(type) {
		case syn.TokenEvent:
			var kind protocol.SemanticTokenType
			switch event.Token.Kind {
			case lex.Let, lex.False, lex.True, lex.Null, lex.Assert, lex.Matches, lex.Fragment, lex.On, lex.Query, lex.Mutation, lex.Set:
				kind = protocol.SemanticTokenTypeKeyword
			case lex.Name:
				scope, ok := scopes.Peek()
				if ok {
					switch scope {
					case ScopeObject:
						kind = protocol.SemanticTokenTypeProperty
					default:
						kind = protocol.SemanticTokenTypeVariable
					}
				} else {
					kind = protocol.SemanticTokenTypeVariable
				}

			case lex.Equals:
				kind = protocol.SemanticTokenTypeOperator
			case lex.Int, lex.Float:
				kind = protocol.SemanticTokenTypeNumber
			case lex.String, lex.BlockString:
				kind = protocol.SemanticTokenTypeString
			default:
				return Highlight{}, false
			}

			return Highlight{Pos: event.Token.Position, TokenKind: kind}, true
		case syn.EnterEvent:
			switch event.Node.(type) {
			case *syn.ObjectExpr, *syn.ObjectPat:
				scopes.Push(ScopeObject)
			}
		case syn.ExitEvent:
			switch event.Node.(type) {
			case *syn.ObjectExpr, *syn.ObjectPat:
				lib.Assert(scopes.MustPop() == ScopeObject)
			}
		}
		return Highlight{}, false
	}).Collect()
}
