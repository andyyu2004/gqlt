package syn

import (
	"fmt"
	"io"

	"andyyu2004/gqlt/lex"

	"github.com/vektah/gqlparser/v2/ast"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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

type MatchesExpr struct {
	Expr Expr
	Pat  Pat
}

var _ Expr = MatchesExpr{}

func (MatchesExpr) isExpr() {}
func (MatchesExpr) isNode() {}

func (expr MatchesExpr) Dump(w io.Writer) {
	expr.Expr.Dump(w)
	io.WriteString(w, " matches ")
	expr.Pat.Dump(w)
}

type ListExpr struct {
	Exprs []Expr
}

var _ Expr = ListExpr{}

func (ListExpr) isExpr() {}
func (ListExpr) isNode() {}

func (expr ListExpr) Dump(w io.Writer) {
	io.WriteString(w, "[")

	for i, expr := range expr.Exprs {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		expr.Dump(w)
	}

	io.WriteString(w, "]")
}

type ObjectExpr struct {
	Fields *orderedmap.OrderedMap[string, Expr]
}

var _ Expr = ObjectExpr{}

func (ObjectExpr) isExpr() {}
func (ObjectExpr) isNode() {}

func (expr ObjectExpr) Dump(w io.Writer) {
	io.WriteString(w, "{")

	for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
		io.WriteString(w, " ")
		io.WriteString(w, entry.Key)
		io.WriteString(w, ": ")
		entry.Value.Dump(w)
	}

	io.WriteString(w, " }")
}

type BinaryExpr struct {
	Op    lex.TokenKind
	Left  Expr
	Right Expr
}

var _ Expr = BinaryExpr{}

func (BinaryExpr) isExpr() {}
func (BinaryExpr) isNode() {}

func (expr BinaryExpr) Dump(w io.Writer) {
	expr.Left.Dump(w)
	io.WriteString(w, " ")
	io.WriteString(w, expr.Op.String())
	io.WriteString(w, " ")
	expr.Right.Dump(w)
}
