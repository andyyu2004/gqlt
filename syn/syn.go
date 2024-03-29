package syn

import (
	"fmt"
	"io"
	"strings"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/movio/gqlt/internal/iterator"
	"github.com/movio/gqlt/internal/lex"
	"github.com/movio/gqlt/internal/slice"
	"github.com/movio/gqlt/memosa/lib"
	"github.com/movio/gqlt/memosa/stack"
)

type Token = lex.Token

type Node interface {
	Child
	isNode()

	Format(io.Writer)
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

func (f File) Format(w io.Writer) {
	for _, stmt := range f.Stmts {
		stmt.Format(w)
		_, _ = io.WriteString(w, ";\n")
	}
}

func (f File) String() string {
	b := strings.Builder{}
	f.Format(&b)
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
			case lexer.Token:
				// it's safe to cast lexer.Token to lex.Token because `lex.Token` is a strict superset of `lexer.Token`
				return TokenEvent{Token{Kind: lex.TokenKind(child.Kind), Value: child.Value, Position: child.Position}}, true
			default:
				panic(fmt.Sprintf("unexpected child type %T", child))
			}
		}

		s, ok := stack.Pop()
		lib.Assert(ok)
		return ExitEvent{s.node}, true
	}
}
