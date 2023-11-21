package syn

import (
	"io"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Pat interface {
	Node
	isPat()
}

type NamePat struct {
	Name string
}

var _ Pat = NamePat{}

func (name NamePat) Dump(w io.Writer) {
	io.WriteString(w, name.Name)
}

func (NamePat) isPat()  {}
func (NamePat) isNode() {}

type ObjectPat struct {
	Fields *orderedmap.OrderedMap[string, Pat]
}

var _ Pat = ObjectPat{}

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

func (ObjectPat) isPat()  {}
func (ObjectPat) isNode() {}
