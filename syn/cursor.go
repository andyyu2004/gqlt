package syn

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/slice"
)

func NewCursor(node Node) *Cursor[Node] {
	return &Cursor[Node]{cur: node, par: nil}
}

type Cursor[T Child] struct {
	par *Cursor[Node]
	cur T
}

func (c *Cursor[T]) Children() []*Cursor[Child] {
	switch cur := any(c.cur).(type) {
	case Node:
		return slice.Map(cur.Children(), func(child Child) *Cursor[Child] {
			return &Cursor[Child]{par: c.cast(), cur: child}
		})
	default:
		return nil
	}
}

func (c *Cursor[T]) Pos() ast.Position {
	return c.cur.Pos()
}

func (c *Cursor[T]) Value() T {
	return c.cur
}

func (c *Cursor[T]) Parent() *Cursor[Node] {
	return c.par
}

// finds the token that contains the given point
func (c *Cursor[T]) TokenAt(point ast.Point) *Cursor[Token] {
	if token, ok := any(c.cur).(Token); ok && c.Pos().Contains(point) {
		return &Cursor[Token]{par: c.par, cur: token}
	}

	for _, child := range c.Children() {
		if child.Pos().Contains(point) {
			return child.TokenAt(point)
		}
	}

	return nil
}

func (c *Cursor[T]) cast() *Cursor[Node] {
	return &Cursor[Node]{par: c.par, cur: any(c.cur).(Node)}
}
