package syn

import (
	"bytes"
	"io"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/iterator"
	"github.com/andyyu2004/gqlt/lex"
	"github.com/andyyu2004/gqlt/slice"
	"github.com/andyyu2004/memosa/lib"
	"github.com/andyyu2004/memosa/stack"
)

func Dump(node Node) string {
	var buf bytes.Buffer
	node.Dump(&buf)
	return buf.String()
}

type Token = lex.Token

type Node interface {
	Child
	isNode()

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

func (File) isNode() {}

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

func Traverse(node Node) iterator.Iterator[Child] {
	type State struct {
		children   Children
		childIndex int
	}
	stack := stack.Stack[*State]{}
	stack.Push(&State{node.Children(), 0})

	return func() (Child, bool) {
		for {
			state, ok := stack.Peek()
			if !ok {
				return nil, false
			}

			for _, child := range state.children[state.childIndex:] {
				state.childIndex++
				switch child := child.(type) {
				case Node:
					stack.Push(&State{child.Children(), 0})
				case lex.Token:
				}
				return child, true
			}

			_, ok = stack.Pop()
			lib.Assert(ok)
		}
	}
}
