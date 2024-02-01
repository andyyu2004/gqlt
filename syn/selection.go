package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
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
	Alias        lexer.Token
	Name         lexer.Token
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

func (f *Field) Children() Children {
	children := Children{f.Alias}
	if f.Name != f.Alias {
		// the parser will make name and alias the same token if no alias is provided
		children = append(children, f.Name)
	}

	return append(children, f.SelectionSet)
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
