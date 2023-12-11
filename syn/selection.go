package syn

import "github.com/andyyu2004/gqlt/gqlparser/ast"

type SelectionSet []Selection

type Selection interface {
	isSelection()
	GetPosition() *ast.Position
}

func (*Field) isSelection()          {}
func (*FragmentSpread) isSelection() {}
func (*InlineFragment) isSelection() {}

func (f *Field) GetPosition() *ast.Position          { return f.Position }
func (s *FragmentSpread) GetPosition() *ast.Position { return s.Position }
func (f *InlineFragment) GetPosition() *ast.Position { return f.Position }

type Field struct {
	Alias        string
	Name         string
	Arguments    ArgumentList
	Directives   DirectiveList
	SelectionSet SelectionSet
	Position     *ast.Position `dump:"-"`
	Comment      *CommentGroup

	// Require validation
	Definition       *FieldDefinition
	ObjectDefinition *Definition
}

type Argument struct {
	Name     string
	Value    *Value
	Position *ast.Position `dump:"-"`
	Comment  *CommentGroup
}

func (f *Field) ArgumentMap(vars map[string]interface{}) map[string]interface{} {
	return arg2map(f.Definition.Arguments, f.Arguments, vars)
}
