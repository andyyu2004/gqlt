// This package is a small wrapper around "github.com/vektah/gqlparser/v2/lexer" with an extended token set
// and a slightly altered interface.
package lex

import (
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/lexer"
	"github.com/vektah/gqlparser/v2/parser"
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

var _ parser.Lexer = new(Lexer)

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

func (l *Lexer) Peek() Token {
	return convertToken(l.tokens[0])
}

func (l *Lexer) Next() Token {
	tok := l.Peek()
	if tok.Kind != EOF {
		// avoid advancing beyond EOF
		l.tokens = l.tokens[1:]
	}
	return tok
}

func convertToken(tok lexer.Token) Token {
	var kind TokenKind
	switch tok.Kind {
	case lexer.Name:
		switch tok.Value {
		case "let":
			kind = Let
		case "query":
			kind = Query
		case "mutation":
			kind = Mutation
		default:
			kind = Name
		}
	default:
		// we have a test asserting the numeric values of these tokens are matching so this conversion is correct
		kind = TokenKind(tok.Kind)
	}
	return Token{Kind: kind, Value: tok.Value, Pos: tok.Pos}
}

type Token struct {
	Kind  TokenKind
	Value string
	Pos   ast.Position
}

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
	Semi
	Equals
	At
	BracketL
	BracketR
	BraceL
	BraceR
	Pipe
	Name
	Int
	Float
	String
	BlockString
	Comment

	// gqlt tokens
	// gqlparser just puts most of these as names, but this is nicer imo
	Let
	Query
	Mutation
)

func (t TokenKind) Name() string {
	switch t {
	case Invalid:
		return "Invalid"
	case EOF:
		return "EOF"
	case Bang:
		return "Bang"
	case Dollar:
		return "Dollar"
	case Amp:
		return "Amp"
	case ParenL:
		return "ParenL"
	case ParenR:
		return "ParenR"
	case Spread:
		return "Spread"
	case Colon:
		return "Colon"
	case Semi:
		return "Semicolon"
	case Equals:
		return "Equals"
	case At:
		return "At"
	case BracketL:
		return "BracketL"
	case BracketR:
		return "BracketR"
	case BraceL:
		return "BraceL"
	case BraceR:
		return "BraceR"
	case Pipe:
		return "Pipe"
	case Name:
		return "Name"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case BlockString:
		return "BlockString"
	case Comment:
		return "Comment"
	case Let:
		return "let"
	case Query:
		return "query"
	case Mutation:
		return "mutation"
	}
	return "Unknown " + strconv.Itoa(int(t))
}

func (t TokenKind) String() string {
	switch t {
	case Invalid:
		return "<Invalid>"
	case EOF:
		return "<EOF>"
	case Bang:
		return "!"
	case Dollar:
		return "$"
	case Amp:
		return "&"
	case ParenL:
		return "("
	case ParenR:
		return ")"
	case Spread:
		return "..."
	case Colon:
		return ":"
	case Semi:
		return ";"
	case Equals:
		return "="
	case At:
		return "@"
	case BracketL:
		return "["
	case BracketR:
		return "]"
	case BraceL:
		return "{"
	case BraceR:
		return "}"
	case Pipe:
		return "|"
	case Name:
		return "Name"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case BlockString:
		return "BlockString"
	case Comment:
		return "Comment"
	case Let:
		return "let"
	case Query:
		return "query"
	case Mutation:
		return "mutation"
	}
	return "Unknown " + strconv.Itoa(int(t))
}

func (t Token) String() string {
	return t.Kind.String()
}
