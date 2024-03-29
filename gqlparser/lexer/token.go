package lexer

import (
	"strconv"

	"github.com/movio/gqlt/gqlparser/ast"
)

const (
	Invalid Type = iota
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
	AngleL
	AngleR
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
)

func (t Type) Name() string {
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
	case Dot:
		return "Dot"
	case Spread:
		return "Spread"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
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
	case AngleL:
		return "AngleL"
	case AngleR:
		return "AngleR"
	case Pipe:
		return "Pipe"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"
	case Tilde:
		return "Tilde"
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
	}
	return "Unknown " + strconv.Itoa(int(t))
}

func (t Type) String() string {
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
	case Dot:
		return "."
	case Spread:
		return "..."
	case Comma:
		return ","
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
	case AngleL:
		return "<"
	case AngleR:
		return ">"
	case Pipe:
		return "|"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Star:
		return "*"
	case Slash:
		return "/"
	case Tilde:
		return "~"
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
	}
	return "Unknown " + strconv.Itoa(int(t))
}

// Kind represents a type of token. The types are predefined as constants.
type Type int

type Token struct {
	Kind     Type         // The token type.
	Value    string       // The literal value consumed.
	Position ast.Position // The file and line this token was read from
}

func (t Token) Dump() string {
	return strconv.Quote(t.Value)
}

func (t Token) Pos() ast.Position {
	return t.Position
}

func (t Token) String() string {
	if t.Value != "" {
		return t.Kind.String() + " " + strconv.Quote(t.Value)
	}
	return t.Kind.String()
}
