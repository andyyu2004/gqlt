package syn

import "io"

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

// implicitly traverse down until there is a node with more than one child
// type DestructurePat {}
