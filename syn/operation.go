package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
	"github.com/andyyu2004/gqlt/lex"
)

type Operation string

const (
	Query        Operation = "query"
	Mutation     Operation = "mutation"
	Subscription Operation = "subscription"
)

type OperationDefinition struct {
	Operation           Operation
	OperationToken      *lex.Token `dump:"-"`
	Name                lexer.Token
	VariableDefinitions VariableDefinitionList
	Directives          DirectiveList
	SelectionSet        SelectionSet
	Position            ast.Position `dump:"-"`
	Comment             *CommentGroup
}

var _ Node = OperationDefinition{}

func (op OperationDefinition) Pos() ast.Position {
	return op.Position
}

func (OperationDefinition) isNode() {}

func (OperationDefinition) Format(io.Writer) {}

func (d OperationDefinition) Children() Children {
	children := Children{}
	if d.OperationToken != nil {
		children = append(children, *d.OperationToken)
	}

	return append(children, d.Name, d.SelectionSet)
}

type VariableDefinition struct {
	Variable     string
	Type         *Type
	DefaultValue *Value
	Directives   DirectiveList
	Position     ast.Position `dump:"-"`
	Comment      *CommentGroup

	// Requires validation
	Definition *Definition
	Used       bool `dump:"-"`
}
