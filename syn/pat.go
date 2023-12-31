package syn

import (
	"fmt"
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/lex"
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
func (WildcardPat) IsNode() {}

func (WildcardPat) Dump(w io.Writer) {
	io.WriteString(w, "_")
}

// A name pattern matches any value and binds it to the name
type NamePat struct {
	ast.Position
	Name lex.Token
}

var _ Pat = NamePat{}

func (pat NamePat) Children() Children {
	return Children{pat.Name}
}

func (NamePat) isPat()  {}
func (NamePat) IsNode() {}

func (pat NamePat) Dump(w io.Writer) {
	io.WriteString(w, pat.Name.Value)
}

// A literal pattern matches a value that is equal to the literal
type LiteralPat struct {
	ast.Position
	Token lex.Token
	Value any
}

var _ Pat = LiteralPat{}

func (LiteralPat) isPat()  {}
func (LiteralPat) IsNode() {}

func (pat LiteralPat) Children() Children {
	return Children{pat.Token}
}

func (name LiteralPat) Dump(w io.Writer) {
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
func (ListPat) IsNode() {}

func (pat ListPat) Dump(w io.Writer) {
	io.WriteString(w, "[")

	for i, pat := range pat.Pats {
		if i > 0 {
			io.WriteString(w, ", ")
		}

		pat.Dump(w)
	}

	io.WriteString(w, "]")
}

type ObjectPat struct {
	ast.Position
	Fields *orderedmap.OrderedMap[lex.Token, Pat]
}

var _ Pat = ObjectPat{}

func (pat ObjectPat) Children() Children {
	children := make(Children, 0, pat.Fields.Len())
	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		children = append(children, entry.Key, entry.Value)
	}
	return children
}

func (ObjectPat) isPat()  {}
func (ObjectPat) IsNode() {}

func (pat ObjectPat) Dump(w io.Writer) {
	io.WriteString(w, "{")

	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		io.WriteString(w, " ")
		io.WriteString(w, entry.Key.Value)
		io.WriteString(w, ": ")
		entry.Value.Dump(w)
	}

	io.WriteString(w, " }")
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
func (RestPat) IsNode() {}

func (pat RestPat) Dump(w io.Writer) {
	io.WriteString(w, "...")
	pat.Pat.Dump(w)
}
