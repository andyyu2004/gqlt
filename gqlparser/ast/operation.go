package ast

import "io"

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
	Position            *Position `dump:"-"`
	Comment             *CommentGroup
}

func (op OperationDefinition) Pos() Position {
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
	Position     *Position `dump:"-"`
	Comment      *CommentGroup

	// Requires validation
	Definition *Definition
	Used       bool `dump:"-"`
}
