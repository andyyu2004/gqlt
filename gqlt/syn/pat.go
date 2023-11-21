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

func (name NamePat) Dump(w io.Writer) error {
	_, err := io.WriteString(w, name.Name)
	return err
}

func (NamePat) isPat()  {}
func (NamePat) isNode() {}

type ObjectPat struct {
	Fields *orderedmap.OrderedMap[string, Pat]
}

var _ Pat = ObjectPat{}

func (pat ObjectPat) Dump(w io.Writer) error {
	if _, err := io.WriteString(w, "{"); err != nil {
		return err
	}

	for entry := pat.Fields.Oldest(); entry != nil; entry = entry.Next() {
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}

		if _, err := io.WriteString(w, entry.Key); err != nil {
			return err
		}

		if _, err := io.WriteString(w, ": "); err != nil {
			return err
		}

		if err := entry.Value.Dump(w); err != nil {
			return err
		}

	}

	_, err := io.WriteString(w, " }")
	return err
}

func (ObjectPat) isPat()  {}
func (ObjectPat) isNode() {}
