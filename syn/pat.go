package syn

import (
	"fmt"
	"io"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/lex"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Pat interface {
	Node
	isPat()
}

type WildcardPat struct {
	ast.Position
}

var _ Pat = WildcardPat{}

func (WildcardPat) Children() Children {
	return nil
}

func (WildcardPat) isPat()  {}
func (WildcardPat) isNode() {}

func (WildcardPat) Format(w io.Writer) {
	_, _ = io.WriteString(w, "_")
}

// A name pattern matches any value and binds it to the name
// Not to be confused with the variable pattern
type NamePat struct {
	Name lex.Token
}

var _ Pat = NamePat{}

func (pat NamePat) Pos() ast.Position {
	return pat.Name.Pos()
}

func (pat NamePat) Children() Children {
	return Children{pat.Name}
}

func (NamePat) isPat()  {}
func (NamePat) isNode() {}

func (pat NamePat) Format(w io.Writer) {
	_, _ = io.WriteString(w, pat.Name.Value)
}

// A variable pattern matches the value of the referenced variable
// Not to be confused with the name pattern
type VarPat struct {
	Dollar lex.Token
	Name   lex.Token
}

var _ Pat = VarPat{}

func (pat VarPat) Pos() ast.Position {
	return pat.Dollar.Pos().Merge(pat.Name.Pos())
}

func (pat VarPat) Children() Children {
	return Children{pat.Dollar, pat.Name}
}

func (pat VarPat) Format(w io.Writer) {
	_, _ = io.WriteString(w, "$")
	_, _ = io.WriteString(w, pat.Name.Value)
}

func (VarPat) isPat()  {}
func (VarPat) isNode() {}

// A literal pattern matches a value that is equal to the literal
type LiteralPat struct {
	ast.Position
	Token lex.Token
	Value any
}

var _ Pat = LiteralPat{}

func (LiteralPat) isPat()  {}
func (LiteralPat) isNode() {}

func (pat LiteralPat) Children() Children {
	return Children{pat.Token}
}

func (name LiteralPat) Format(w io.Writer) {
	fmt.Fprintf(w, "%v", name.Value)
}

type ListPat struct {
	ast.Position
	Pats []Pat
}

var _ Pat = ListPat{}

func (pat ListPat) Children() Children {
	children := make(Children, len(pat.Pats))
	for i, pat := range pat.Pats {
		children[i] = pat
	}
	return children
}

func (ListPat) isPat()  {}
func (ListPat) isNode() {}

func (pat ListPat) Format(w io.Writer) {
	_, _ = io.WriteString(w, "[")

	for i, pat := range pat.Pats {
		if i > 0 {
			_, _ = io.WriteString(w, ", ")
		}

		pat.Format(w)
	}

	_, _ = io.WriteString(w, "]")
}

type ObjectPat struct {
	ast.Position
	Fields *orderedmap.OrderedMap[lex.Token, Pat]
}

var _ Pat = ObjectPat{}

func (pat ObjectPat) Children() Children {
	children := make(Children, 0, pat.Fields.Len())
	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		// Only include the key in here as the pattern's position does not overlap with it
		// This is an edge case arising from punning
		if entry.Key.Position != entry.Value.Pos() {
			children = append(children, entry.Key)
		}
		children = append(children, entry.Value)
	}
	return children
}

func (ObjectPat) isPat()  {}
func (ObjectPat) isNode() {}

func (pat ObjectPat) Format(w io.Writer) {
	_, _ = io.WriteString(w, "{")

	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		_, _ = io.WriteString(w, " ")
		_, _ = io.WriteString(w, entry.Key.Value)
		_, _ = io.WriteString(w, ": ")
		entry.Value.Format(w)
	}

	_, _ = io.WriteString(w, " }")
}

type RestPat struct {
	ast.Position
	Pat Pat
}

var _ Pat = RestPat{}

func (pat RestPat) Children() Children {
	return Children{pat.Pat}
}

func (RestPat) isPat()  {}
func (RestPat) isNode() {}

func (pat RestPat) Format(w io.Writer) {
	_, _ = io.WriteString(w, "...")
	pat.Pat.Format(w)
}
