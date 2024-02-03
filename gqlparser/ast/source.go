package ast

import "fmt"

// Source covers a single *.graphql file
type Source struct {
	// Name is the filename of the source
	Name string
	// Input is the actual contents of the source file
	Input string
	// BuiltIn indicate whether the source is a part of the specification
	BuiltIn bool
}

type Point = int

type Position struct {
	Start  int     // The starting position, in runes, of this token in the input.
	End    int     // The end position, in runes, of this token in the input.
	Line   int     // The line number at the start of this item.
	Column int     // The column number at the start of this item.
	Src    *Source // The source document this token belongs to
}

func (p Position) Contains(point Point) bool {
	return p.Start <= point && point < p.End
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Src.Name, p.Line, p.Column)
}

func (p Position) Pos() Position {
	return p
}

type HasPosition interface {
	Pos() Position
}

func (p Position) Merge(other HasPosition) Position {
	pos := other.Pos()

	if p.Start > pos.Start {
		p.Start = pos.Start
	}

	if p.End < pos.End {
		p.End = pos.End
	}

	return p
}
