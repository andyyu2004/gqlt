package syn

import "io"

type Node interface {
	isNode()

	Dump(io.Writer)
}

type File struct {
	Stmts []Stmt
}

func (f File) Dump(w io.Writer) {
	for i, stmt := range f.Stmts {
		if i > 0 {
			io.WriteString(w, "\n")
		}

		stmt.Dump(w)
	}
}
