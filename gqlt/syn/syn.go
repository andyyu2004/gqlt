package syn

import "io"

type Node interface {
	isNode()

	Dump(io.Writer) error
}

type File struct {
	Stmts []Stmt
}

func (f File) Dump(w io.Writer) error {
	for i, stmt := range f.Stmts {
		if i > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}

		if err := stmt.Dump(w); err != nil {
			return err
		}

	}
	return nil
}
