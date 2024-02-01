package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

type SelectionSet []Selection

var _ Node = SelectionSet{}

func (ss SelectionSet) Dump(w io.Writer) {
	for _, s := range ss {
		s.Dump(w)
		_, _ = io.WriteString(w, "\n")
	}
}

func (ss SelectionSet) Pos() ast.Position {
	fst := ss[0]
	lst := ss[len(ss)-1]
	return fst.Pos().Merge(lst)
}

func (SelectionSet) isNode() {}

func (s SelectionSet) Children() Children {
	children := make(Children, 0, len(s))
	for _, sel := range s {
		children = append(children, sel)
	}
	return children
}

type Selection interface {
	Node
	isSelection()
}

func (*Field) isSelection()          {}
func (*FragmentSpread) isSelection() {}
func (*InlineFragment) isSelection() {}

func (f *Field) Pos() ast.Position          { return f.Position }
func (s *FragmentSpread) Pos() ast.Position { return s.Position }
func (f *InlineFragment) Pos() ast.Position { return f.Position }

type Field struct {
	Alias        string
	Name         string
	Arguments    ArgumentList
	Directives   DirectiveList
	SelectionSet SelectionSet
	Position     ast.Position `dump:"-"`
	Comment      *CommentGroup

	// Require validation
	Definition       *FieldDefinition
	ObjectDefinition *Definition
}

var _ Node = &Field{}

func (*Field) Children() Children {
	return Children{}
}

func (*Field) Dump(io.Writer) {
}

func (*Field) isNode() {}

type Argument struct {
	Name     string
	Value    *Value
	Position ast.Position `dump:"-"`
	Comment  *CommentGroup
}

func (f *Field) ArgumentMap(vars map[string]interface{}) map[string]interface{} {
	return arg2map(f.Definition.Arguments, f.Arguments, vars)
}
