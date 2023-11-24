package syn

import (
	"fmt"
	"io"

	"github.com/andyyu2004/gqlt/lex"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Expr interface {
	Node
	isExpr()
}

type NameExpr struct {
	ast.Position
	Name string
}

var _ Expr = NameExpr{}

func (NameExpr) isExpr() {}
func (NameExpr) isNode() {}

func (expr NameExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Name)
}

type OperationExpr struct {
	ast.Position
	// unparsed graphql string
	// useful for pretty printing without formatting
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

type IndexExpr struct {
	ast.Position
	Expr  Expr
	Index Expr
}

var _ Expr = IndexExpr{}

func (IndexExpr) isExpr() {}
func (IndexExpr) isNode() {}

func (expr IndexExpr) Dump(w io.Writer) {
	expr.Expr.Dump(w)
	io.WriteString(w, "[")
	expr.Index.Dump(w)
	io.WriteString(w, "]")
}

type LiteralExpr struct {
	ast.Position
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
	ast.Position
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
	ast.Position
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
	ast.Position
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
	ast.Position
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
	ast.Position
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

type UnaryExpr struct {
	ast.Position
	Op   lex.TokenKind
	Expr Expr
}

var _ Expr = UnaryExpr{}

func (UnaryExpr) isExpr() {}
func (UnaryExpr) isNode() {}

func (expr UnaryExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Op.String())
	expr.Expr.Dump(w)
}
