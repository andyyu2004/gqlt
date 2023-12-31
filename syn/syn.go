package syn

import (
	"io"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/iterator"
	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/slice"
	"github.com/andyyu2004/memosa/lib"
	"github.com/andyyu2004/memosa/stack"
)

// func Dump(node Node) string {
// 	var buf bytes.Buffer
// 	node.Dump(&buf)
// 	return buf.String()
// }

type Token = lex.Token

type Node interface {
	Child
	IsNode()

	Dump(io.Writer)
	Children() Children
}

type File struct {
	ast.Position
	Stmts []Stmt
}

func (f File) Children() Children {
	return slice.Map(f.Stmts, func(stmt Stmt) Child { return stmt })
}

func (File) IsNode() {}

var _ Node = File{}

type Children []Child

// Node | Token
type Child interface {
	ast.HasPosition
}

func (f File) Dump(w io.Writer) {
	for _, stmt := range f.Stmts {
		stmt.Dump(w)
		io.WriteString(w, ";\n")
	}
}

func (f File) String() string {
	b := strings.Builder{}
	f.Dump(&b)
	return b.String()
}

type Event interface {
	isEvent()
}

type EnterEvent struct{ Node Node }

func (EnterEvent) isEvent() {}

type ExitEvent struct{ Node Node }

func (ExitEvent) isEvent() {}

type TokenEvent struct{ Token Token }

func (TokenEvent) isEvent() {}

func Traverse(node Node) iterator.Iterator[Event] {
	type State struct {
		node       Node
		children   Children
		childIndex int
	}
	stack := stack.Stack[*State]{}
	stack.Push(&State{node, node.Children(), 0})

	return func() (Event, bool) {
		state, ok := stack.Peek()
		if !ok {
			return nil, false
		}

		for _, child := range state.children[state.childIndex:] {
			state.childIndex++
			switch child := child.(type) {
			case Node:
				stack.Push(&State{child, child.Children(), 0})
				return EnterEvent{child}, true
			case lex.Token:
				return TokenEvent{child}, true
			default:
				panic("unreachable")
			}
		}

		s, ok := stack.Pop()
		lib.Assert(ok)
		return ExitEvent{s.node}, true
	}
}
