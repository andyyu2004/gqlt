package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/lex"
)

func NonNullNamedType(named string, pos ast.Position) *Type {
	return &Type{NamedType: lex.Token{Value: named, Kind: lex.TypeName}, NonNull: true, Position: pos}
}

func NamedType(named string, pos ast.Position) *Type {
	return &Type{NamedType: lex.Token{Value: named, Kind: lex.TypeName}, NonNull: false, Position: pos}
}

func NonNullListType(elem *Type, pos ast.Position) *Type {
	return &Type{Elem: elem, NonNull: true, Position: pos}
}

func ListType(elem *Type, pos ast.Position) *Type {
	return &Type{Elem: elem, NonNull: false, Position: pos}
}

type Type struct {
	NamedType lex.Token
	Elem      *Type
	NonNull   bool
	Position  ast.Position `dump:"-"`
}

var _ Node = Type{}

func (t Type) Children() Children {
	children := Children{}
	if t.NamedType.Value != "" {
		children = append(children, t.NamedType)
	}
	if t.Elem != nil {
		children = append(children, t.Elem)
	}

	return children
}

func (Type) Format(io.Writer) {}

func (t Type) Pos() ast.Position {
	return t.Position
}

func (Type) isNode() {}

func (t *Type) Name() string {
	if t.NamedType.Value != "" {
		return t.NamedType.Value
	}

	return t.Elem.Name()
}

func (t *Type) String() string {
	nn := ""
	if t.NonNull {
		nn = "!"
	}
	if t.NamedType.Value != "" {
		return t.NamedType.Value + nn
	}

	return "[" + t.Elem.String() + "]" + nn
}

func (t *Type) IsCompatible(other *Type) bool {
	if t.NamedType.Value != other.NamedType.Value {
		return false
	}

	if t.Elem != nil && other.Elem == nil {
		return false
	}

	if t.Elem != nil && !t.Elem.IsCompatible(other.Elem) {
		return false
	}

	if other.NonNull {
		return t.NonNull
	}

	return true
}

func (t *Type) Dump() string {
	return t.String()
}
