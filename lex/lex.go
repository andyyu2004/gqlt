// This package is a small wrapper around "github.com/andyyu2004/gqlt/gqlparser/lexer" with an extended token set
// and a slightly altered interface.
package lex

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
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
	lex := lexer.New(src)
	for {
		tok, err := lex.ReadToken()
		if err != nil {
			return Lexer{}, err
		}

		if tok.Kind != lexer.Comment {
			tokens = append(tokens, tok)
		}

		if tok.Kind == lexer.EOF {
			break
		}
	}

	assert(len(tokens) > 0)

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
		case "let":
			kind = Let
		case "query":
			kind = Query
		case "mutation":
			kind = Mutation
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
		default:
			kind = Name
		}
	case lexer.Equals:
		if tok[1].Kind == lexer.Equals {
			kind = Equals2
			len = 2
		}
	case lexer.Bang:
		if tok[1].Kind == lexer.Equals {
			kind = BangEqual
			len = 2
		}

	}
	return Token{Kind: kind, Value: tok[0].Value, Position: tok[0].Pos}, len
}

type Token struct {
	Kind  TokenKind
	Value string
	ast.Position
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
	Name
	Int
	Float
	String
	BlockString
	Comment

	// gqlt tokens
	// gqlparser just puts most of these as names, but this is nicer imo
	Underscore
	Let
	Query
	Mutation
	True
	False
	Null
	Matches
	Assert
	Set
	Equals2
	BangEqual
	Not
)

func (t TokenKind) Name() string {
	switch t {
	case Underscore:
		return "_"
	case Let:
		return "let"
	case Query:
		return "query"
	case Mutation:
		return "mutation"
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
	case Query:
		return "query"
	case Mutation:
		return "mutation"
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
	case Not:
		return "not"
	default:
		return lexer.Type(t).String()
	}
}

func (t Token) String() string {
	return t.Kind.String()
}
