// This package is a small wrapper around "github.com/movio/gqlt/gqlparser/lexer" with an extended token set
// and a slightly altered interface.
package lex

import (
	"strconv"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/gqlerror"
	"github.com/movio/gqlt/gqlparser/lexer"
)

type Lexer struct {
	tokens []lexer.Token
}

// compatibility with gqlparser's lexer
func (l *Lexer) ReadToken() (lexer.Token, error) {
	tok := l.tokens[0]
	if tok.Kind != lexer.EOF {
		l.tokens = l.tokens[1:]
	}

	return tok, nil
}

func New(src *ast.Source) (Lexer, error) {
	// Read all the tokens eagerly to save us error handling on `ReadToken` later.
	// We're not expecting large files so this should be fine.
	var tokens []lexer.Token
	var errors ast.Errors
	lex := lexer.New(src)
	for {
		tok, err := lex.ReadToken()
		if err != nil {
			err := err.(*gqlerror.Error)
			errors = append(errors, ast.Error{Position: tok.Position, Msg: err.Message})
			continue
		}

		if tok.Kind != lexer.Comment {
			tokens = append(tokens, tok)
		}

		if tok.Kind == lexer.EOF {
			break
		}
	}

	assert(len(tokens) > 0)

	if len(errors) > 0 {
		return Lexer{tokens: tokens}, errors
	}

	return Lexer{tokens: tokens}, nil
}

func assert(ok bool) {
	if !ok {
		panic("lexer assertion failed")
	}
}

const n = 2

func (l *Lexer) Peek() Token {
	token, _ := l.peek()
	return token
}

func (l *Lexer) peek() (Token, int) {
	var ts [n]lexer.Token
	if len(l.tokens) > 1 {
		ts = [n]lexer.Token(l.tokens[:n])
	} else {
		ts = [n]lexer.Token{l.tokens[0], {Kind: lexer.EOF}}
	}
	return convertToken(ts)
}

func (l *Lexer) Next() Token {
	tok, len := l.peek()
	if tok.Kind != EOF {
		// avoid advancing beyond EOF
		l.tokens = l.tokens[len:]
	}
	return tok
}

// we take 2 tokens in as some gqlt tokens are composed of up to `n` tokens
func convertToken(tok [n]lexer.Token) (Token, int) {
	// we have a test asserting the numeric values of these tokens are matching so this conversion is correct
	kind := TokenKind(tok[0].Kind)
	len := 1
	switch tok[0].Kind {
	case lexer.Name:
		switch tok[0].Value {
		case "_":
			kind = Underscore
		case ".":
			kind = Dot
		case "let":
			kind = Let
		case "use":
			kind = Use
		case "query":
			kind = Query
		case "mutation":
			kind = Mutation
		case "fragment":
			kind = Fragment
		case "on":
			kind = On
		case "true":
			kind = True
		case "false":
			kind = False
		case "null":
			kind = Null
		case "matches":
			kind = Matches
		case "assert":
			kind = Assert
		case "set":
			kind = Set
		case "not":
			kind = Not
		case "try":
			kind = Try
		default:
			kind = Name
		}
	case lexer.Equals:
		switch tok[1].Kind {
		case lexer.Equals:
			kind = Equals2
			len = 2
		case lexer.Tilde:
			kind = EqualsTilde
			len = 2
		}
	case lexer.Bang:
		switch tok[1].Kind {
		case lexer.Equals:
			kind = BangEqual
			len = 2
		case lexer.Tilde:
			kind = BangTilde
			len = 2
		}
	}
	return Token{Kind: kind, Value: tok[0].Value, Position: tok[0].Position}, len
}

type Token struct {
	Kind  TokenKind
	Value string
	ast.Position
}

func (t Token) Dump() string {
	return strconv.Quote(t.Value)
}

func (Token) IsChild() {}

type TokenKind lexer.Type

const (
	Invalid TokenKind = iota
	EOF
	Bang
	Dollar
	Amp
	ParenL
	ParenR
	Dot
	Spread
	Colon
	Comma
	Semi
	Equals
	At
	BracketL
	BracketR
	BraceL
	BraceR
	Pipe
	Plus
	Minus
	Star
	Slash
	Tilde
	Name
	Int
	Float
	String
	BlockString
	Comment
	// gqlparser tokens: every token above here must be consistent with `lexer.Token` so they can be safely cast

	// gqlt tokens
	// gqlparser just puts most of these as names, but this is nicer imo
	Underscore
	Let
	Use
	Query
	Mutation
	Subscription
	Fragment
	On
	True
	False
	Null
	Matches
	Assert
	Set
	Equals2
	BangEqual
	Not
	Try
	EqualsTilde
	BangTilde

	TypeName
)

func (t TokenKind) Name() string {
	switch t {
	case Underscore:
		return "_"
	case Let:
		return "let"
	case Use:
		return "use"
	case Query:
		return "query"
	case Mutation:
		return "mutation"
	case Subscription:
		return "subscription"
	case Fragment:
		return "fragment"
	case On:
		return "on"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"
	case Matches:
		return "matches"
	case Assert:
		return "assert"
	case Set:
		return "set"
	case Equals2:
		return "Equals2"
	case BangEqual:
		return "BangEqual"
	case Not:
		return "Not"
	case Try:
		return "try"
	default:
		return lexer.Type(t).Name()
	}
}

func (t TokenKind) String() string {
	switch t {
	case Underscore:
		return "_"
	case Let:
		return "let"
	case Use:
		return "use"
	case Query:
		return "query"
	case Mutation:
		return "mutation"
	case Subscription:
		return "subscription"
	case Fragment:
		return "fragment"
	case On:
		return "on"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "null"
	case Matches:
		return "matches"
	case Assert:
		return "assert"
	case Set:
		return "set"
	case Equals2:
		return "=="
	case BangEqual:
		return "!="
	case EqualsTilde:
		return "=~"
	case BangTilde:
		return "!~"
	case Not:
		return "not"
	case Try:
		return "try"
	default:
		return lexer.Type(t).String()
	}
}

func (t Token) String() string {
	return t.Kind.String()
}
