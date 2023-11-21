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

type NameExpr struct {
	Name string
}

var _ Expr = NameExpr{}

func (NameExpr) isExpr() {}
func (NameExpr) isNode() {}

func (expr NameExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Name)
}

type OperationExpr struct {
	// unparsed graphql string
	Query string
	// parsed graphql ast
	Operation *ast.OperationDefinition
}

var _ Expr = OperationExpr{}

func (OperationExpr) isExpr() {}
func (OperationExpr) isNode() {}

func (expr OperationExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Query)
}

type LiteralExpr struct {
	Value any
}

var _ Expr = LiteralExpr{}

func (LiteralExpr) isExpr() {}
func (LiteralExpr) isNode() {}

func (expr LiteralExpr) Dump(w io.Writer) {
	switch expr.Value.(type) {
	case string:
		fmt.Fprintf(w, "\"%v\"", expr.Value)
	default:
		fmt.Fprintf(w, "%v", expr.Value)
	}
}

type CallExpr struct {
	Fn   Expr
	Args []Expr
}

var _ Expr = CallExpr{}

func (CallExpr) isExpr() {}
func (CallExpr) isNode() {}

func (expr CallExpr) Dump(w io.Writer) {
	expr.Fn.Dump(w)
	io.WriteString(w, "(")

	for i, arg := range expr.Args {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		arg.Dump(w)
	}
	io.WriteString(w, ")")
}
