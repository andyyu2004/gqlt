package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/lex"
)

type FragmentSpread struct {
	Name       string
	Directives DirectiveList

	// Require validation
	ObjectDefinition *Definition
	Definition       *FragmentDefinition

	Position ast.Position `dump:"-"`
	Comment  *CommentGroup
}

var _ Node = &FragmentSpread{}

func (FragmentSpread) Children() Children {
	return Children{}
}

func (FragmentSpread) Dump(io.Writer) {
}

func (FragmentSpread) isNode() {}

type InlineFragment struct {
	TypeCondition string
	Directives    DirectiveList
	SelectionSet  SelectionSet

	// Require validation
	ObjectDefinition *Definition

	Position ast.Position `dump:"-"`
	Comment  *CommentGroup
}

var _ Node = &InlineFragment{}

func (InlineFragment) Children() Children {
	return Children{}
}

func (InlineFragment) Dump(io.Writer) {
}

func (InlineFragment) isNode() {}

type FragmentDefinition struct {
	Name lex.Token
	// Note: fragment variable definitions are experimental and may be changed
	// or removed in the future.
	VariableDefinition VariableDefinitionList
	TypeCondition      lex.Token
	Directives         DirectiveList
	SelectionSet       SelectionSet

	// Require validation
	Definition *Definition

	FragmentKw Token        `dump:"-"`
	OnKw       Token        `dump:"-"`
	Position   ast.Position `dump:"-"`
	Comment    *CommentGroup
}

var _ Node = FragmentDefinition{}

func (f FragmentDefinition) Children() Children {
	return Children{
		f.FragmentKw,
		f.Name,
		f.OnKw,
		f.TypeCondition,
		f.SelectionSet,
	}
}

func (d FragmentDefinition) Dump(w io.Writer) {
	_, _ = io.WriteString(w, "fragment ")
	_, _ = io.WriteString(w, d.Name.Value)
	_, _ = io.WriteString(w, "on ")
	_, _ = io.WriteString(w, d.TypeCondition.Value)
	_, _ = io.WriteString(w, "{ ... }")
}

func (f FragmentDefinition) Pos() ast.Position {
	return f.Position
}

func (FragmentDefinition) isNode() {}
