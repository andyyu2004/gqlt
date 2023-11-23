package syn

import (
	"fmt"
	"io"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Pat interface {
	Node
	isPat()
}

type WildcardPat struct{}

var _ Pat = WildcardPat{}

func (WildcardPat) isPat()  {}
func (WildcardPat) isNode() {}

func (WildcardPat) Dump(w io.Writer) {
	io.WriteString(w, "_")
}

// A name pattern matches any value and binds it to the name
type NamePat struct {
	Name string
}

var _ Pat = NamePat{}

func (NamePat) isPat()  {}
func (NamePat) isNode() {}

func (name NamePat) Dump(w io.Writer) {
	io.WriteString(w, name.Name)
}

// A literal pattern matches a value that is equal to the literal
type LiteralPat struct {
	Value any
}

var _ Pat = LiteralPat{}

func (LiteralPat) isPat()  {}
func (LiteralPat) isNode() {}

func (name LiteralPat) Dump(w io.Writer) {
	fmt.Fprintf(w, "%v", name.Value)
}

type ListPat struct {
	Pats []Pat
}

var _ Pat = ListPat{}

func (ListPat) isPat()  {}
func (ListPat) isNode() {}

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
	Fields *orderedmap.OrderedMap[string, Pat]
}

var _ Pat = ObjectPat{}

func (ObjectPat) isPat()  {}
func (ObjectPat) isNode() {}

func (pat ObjectPat) Dump(w io.Writer) {
	io.WriteString(w, "{")

	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		io.WriteString(w, " ")
		io.WriteString(w, entry.Key)
		io.WriteString(w, ": ")
		entry.Value.Dump(w)
	}

	io.WriteString(w, " }")
}
