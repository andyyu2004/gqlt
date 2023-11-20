package syn

import (
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
