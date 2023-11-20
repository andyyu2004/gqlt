package syn

import "io"

type Stmt interface {
	Node
	isStmt()
}

type LetStmt struct {
	Pat  Pat
	Expr Expr
}

var _ Stmt = LetStmt{}

func (LetStmt) isStmt() {}
func (LetStmt) isNode() {}
func (let LetStmt) Dump(w io.Writer) error {
	if _, err := io.WriteString(w, "let "); err != nil {
		return err
	}

	if err := let.Pat.Dump(w); err != nil {
		return err
	}

	if _, err := io.WriteString(w, " = "); err != nil {
		return err
	}

	if err := let.Expr.Dump(w); err != nil {
		return err
	}

	return nil
}
