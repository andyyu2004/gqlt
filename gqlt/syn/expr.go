package syn

import (
	"fmt"
	"io"

	"github.com/vektah/gqlparser/v2/ast"
)

type Expr interface {
	Node
	isExpr()
}

type OperationExpr struct {
	// unparsed graphql string
	Query string
	// parsed graphql ast
	Operation *ast.OperationDefinition
}

var _ Expr = OperationExpr{}

func (expr OperationExpr) Dump(w io.Writer) error {
	_, err := io.WriteString(w, expr.Query)
	return err
}

func (OperationExpr) isExpr() {}
func (OperationExpr) isNode() {}

type LiteralExpr struct {
	Value any
}

var _ Expr = LiteralExpr{}

func (expr LiteralExpr) Dump(w io.Writer) (err error) {
	switch expr.Value.(type) {
	case string:
		_, err = fmt.Fprintf(w, "\"%v\"", expr.Value)
	default:
		_, err = fmt.Fprintf(w, "%v", expr.Value)
	}
	return
}

func (LiteralExpr) isExpr() {}
func (LiteralExpr) isNode() {}
