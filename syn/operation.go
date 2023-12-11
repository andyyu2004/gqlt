package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

type Operation string

const (
	Query        Operation = "query"
	Mutation     Operation = "mutation"
	Subscription Operation = "subscription"
)

type OperationDefinition struct {
	Operation           Operation
	Name                string
	VariableDefinitions VariableDefinitionList
	Directives          DirectiveList
	SelectionSet        SelectionSet
	Position            *ast.Position `dump:"-"`
	Comment             *CommentGroup
}

func (op OperationDefinition) Pos() ast.Position {
	return *op.Position
}

func (OperationDefinition) IsNode() {}

func (OperationDefinition) Dump(io.Writer) {}

func (OperationDefinition) Children() {}

type VariableDefinition struct {
	Variable     string
	Type         *Type
	DefaultValue *Value
	Directives   DirectiveList
	Position     *ast.Position `dump:"-"`
	Comment      *CommentGroup

	// Requires validation
	Definition *Definition
	Used       bool `dump:"-"`
}
