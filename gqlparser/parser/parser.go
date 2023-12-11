package parser

import (
	"strconv"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/gqlerror"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
	"github.com/andyyu2004/gqlt/syn"
)

// Commas exist in the lexer for use by the gqlt parser.
// However, graphql treats them as whitespace so we need to ignore it.
type skipCommas struct {
	lexer Lexer
}

func (s skipCommas) ReadToken() (lexer.Token, error) {
	for {
		tok, err := s.lexer.ReadToken()
		if err != nil {
			return tok, err
		}

		if tok.Kind != lexer.Comma {
			return tok, err
		}
	}
}

func New(lexer Lexer) *parser {
	return &parser{lexer: skipCommas{lexer}}
}

type Lexer interface {
	ReadToken() (lexer.Token, error)
}

type parser struct {
	lexer Lexer
	err   error

	peeked    bool
	peekToken lexer.Token
	peekError error

	prev lexer.Token

	comment          *syn.CommentGroup
	commentConsuming bool
}

func (p *parser) Err() error {
	return p.err
}

func (p *parser) consumeComment() (*syn.Comment, bool) {
	if p.err != nil {
		return nil, false
	}
	tok := p.peek()
	if tok.Kind != lexer.Comment {
		return nil, false
	}
	p.next()
	return &syn.Comment{
		Value:    tok.Value,
		Position: &tok.Pos,
	}, true
}

func (p *parser) consumeCommentGroup() {
	if p.err != nil {
		return
	}
	if p.commentConsuming {
		return
	}
	p.commentConsuming = true

	var comments []*syn.Comment
	for {
		comment, ok := p.consumeComment()
		if !ok {
			break
		}
		comments = append(comments, comment)
	}

	p.comment = &syn.CommentGroup{List: comments}
	p.commentConsuming = false
}

func (p *parser) peekPos() *ast.Position {
	if p.err != nil {
		return nil
	}

	peek := p.peek()
	return &peek.Pos
}

func (p *parser) peek() lexer.Token {
	if p.err != nil {
		return p.prev
	}

	if !p.peeked {
		p.peekToken, p.peekError = p.lexer.ReadToken()
		p.peeked = true
		if p.peekToken.Kind == lexer.Comment {
			p.consumeCommentGroup()
		}
	}

	return p.peekToken
}

func (p *parser) error(tok lexer.Token, format string, args ...interface{}) {
	if p.err != nil {
		return
	}
	p.err = gqlerror.ErrorLocf(tok.Pos.Src.Name, tok.Pos.Line, tok.Pos.Column, format, args...)
}

func (p *parser) next() lexer.Token {
	if p.err != nil {
		return p.prev
	}
	if p.peeked {
		p.peeked = false
		p.comment = nil
		p.prev, p.err = p.peekToken, p.peekError
	} else {
		p.prev, p.err = p.lexer.ReadToken()
		if p.prev.Kind == lexer.Comment {
			p.consumeCommentGroup()
		}
	}
	return p.prev
}

func (p *parser) expectKeyword(value string) (lexer.Token, *syn.CommentGroup) {
	tok := p.peek()
	comment := p.comment
	if tok.Kind == lexer.Name && tok.Value == value {
		return p.next(), comment
	}

	p.error(tok, "Expected %s, found %s", strconv.Quote(value), tok.String())
	return tok, comment
}

func (p *parser) expect(kind lexer.Type) (lexer.Token, *syn.CommentGroup) {
	tok := p.peek()
	comment := p.comment
	if tok.Kind == kind {
		return p.next(), comment
	}

	p.error(tok, "Expected %s, found %s", kind, tok.Kind.String())
	return tok, comment
}

func (p *parser) skip(kind lexer.Type) bool {
	if p.err != nil {
		return false
	}

	tok := p.peek()

	if tok.Kind != kind {
		return false
	}
	p.next()
	return true
}

func (p *parser) unexpectedError() {
	p.unexpectedToken(p.peek())
}

func (p *parser) unexpectedToken(tok lexer.Token) {
	p.error(tok, "Unexpected %s", tok.String())
}

func (p *parser) many(start lexer.Type, end lexer.Type, cb func()) {
	hasDef := p.skip(start)
	if !hasDef {
		return
	}

	for p.peek().Kind != end && p.err == nil {
		cb()
	}
	p.next()
}

func (p *parser) some(start lexer.Type, end lexer.Type, cb func()) *syn.CommentGroup {
	hasDef := p.skip(start)
	if !hasDef {
		return nil
	}

	called := false
	for p.peek().Kind != end && p.err == nil {
		called = true
		cb()
	}

	if !called {
		p.error(p.peek(), "expected at least one definition, found %s", p.peek().Kind.String())
		return nil
	}

	comment := p.comment
	p.next()
	return comment
}
