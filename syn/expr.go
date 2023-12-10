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
	Name lex.Token
}

var _ Expr = NameExpr{}

func (expr NameExpr) Children() Children {
	return Children{expr.Name}
}

func (NameExpr) isExpr() {}
func (NameExpr) isNode() {}

func (expr NameExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Name.Value)
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

func (expr OperationExpr) Children() Children {
	return Children{
		lex.Token{Kind: lex.Query, Value: expr.Query, Position: expr.Position},
	}
}

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

func (expr IndexExpr) Children() Children {
	return Children{
		expr.Expr,
		expr.Index,
	}
}

func (IndexExpr) isExpr() {}
func (IndexExpr) isNode() {}

func (expr IndexExpr) Dump(w io.Writer) {
	expr.Expr.Dump(w)
	io.WriteString(w, "[")
	expr.Index.Dump(w)
	io.WriteString(w, "]")
}

type LiteralExpr struct {
	lex.Token
	Value any
}

var _ Expr = LiteralExpr{}

func (expr LiteralExpr) Children() Children {
	return Children{expr.Token}
}

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

func (expr CallExpr) Children() Children {
	children := make(Children, 0, len(expr.Args)+1)
	children = append(children, expr.Fn)
	for _, arg := range expr.Args {
		children = append(children, arg)
	}
	return children
}

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

func (expr MatchesExpr) Children() Children {
	return Children{expr.Expr, expr.Pat}
}

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

func (expr ListExpr) Children() Children {
	children := make(Children, len(expr.Exprs))
	for i, expr := range expr.Exprs {
		children[i] = expr
	}
	return children
}

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
	Fields *orderedmap.OrderedMap[lex.Token, Expr]
}

var _ Expr = ObjectExpr{}

func (expr ObjectExpr) Children() Children {
	children := make(Children, 0, expr.Fields.Len()*2)
	for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
		children = append(children, entry.Key, entry.Value)
	}
	return children
}

func (ObjectExpr) isExpr() {}
func (ObjectExpr) isNode() {}

func (expr ObjectExpr) Dump(w io.Writer) {
	io.WriteString(w, "{")

	for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
		io.WriteString(w, " ")
		io.WriteString(w, entry.Key.Value)
		io.WriteString(w, ": ")
		entry.Value.Dump(w)
	}

	io.WriteString(w, " }")
}

type BinaryExpr struct {
	ast.Position
	Op    lex.Token
	Left  Expr
	Right Expr
}

var _ Expr = BinaryExpr{}

func (expr BinaryExpr) Children() Children {
	return Children{expr.Left, expr.Op, expr.Right}
}

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
	Op   lex.Token
	Expr Expr
}

var _ Expr = UnaryExpr{}

func (expr UnaryExpr) Children() Children {
	return Children{expr.Op, expr.Expr}
}

func (UnaryExpr) isExpr() {}
func (UnaryExpr) isNode() {}

func (expr UnaryExpr) Dump(w io.Writer) {
	io.WriteString(w, expr.Op.String())
	expr.Expr.Dump(w)
}
