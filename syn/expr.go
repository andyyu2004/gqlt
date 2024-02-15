package syn

import (
	"fmt"
	"io"

	"github.com/movio/gqlt/internal/lex"

	"github.com/movio/gqlt/gqlparser/ast"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Expr interface {
	Node
	isExpr()
}

type NameExpr struct {
	Name lex.Token
}

var _ Expr = NameExpr{}

func (expr NameExpr) Pos() ast.Position {
	return expr.Name.Pos()
}

func (expr NameExpr) Children() Children {
	return Children{expr.Name}
}

func (NameExpr) isExpr() {}
func (NameExpr) isNode() {}

func (expr NameExpr) Format(w io.Writer) {
	_, _ = io.WriteString(w, expr.Name.Value)
}

type QueryExpr struct {
	ast.Position
	// unparsed graphql string
	// useful for pretty printing without formatting
	Query string
	// parsed graphql ast
	Operation *OperationDefinition
}

var _ Expr = QueryExpr{}

func (expr QueryExpr) Children() Children {
	return Children{expr.Operation}
}

func (QueryExpr) isExpr() {}
func (QueryExpr) isNode() {}

func (expr QueryExpr) Format(w io.Writer) {
	_, _ = io.WriteString(w, expr.Query)
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

func (expr IndexExpr) Format(w io.Writer) {
	expr.Expr.Format(w)
	_, _ = io.WriteString(w, "[")
	expr.Index.Format(w)
	_, _ = io.WriteString(w, "]")
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

func (expr LiteralExpr) Format(w io.Writer) {
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

func (expr CallExpr) Format(w io.Writer) {
	expr.Fn.Format(w)
	_, _ = io.WriteString(w, "(")

	for i, arg := range expr.Args {
		if i > 0 {
			_, _ = io.WriteString(w, ", ")
		}
		arg.Format(w)
	}
	_, _ = io.WriteString(w, ")")
}

type FieldExpr struct {
	Expr  Expr
	Dot   lex.Token
	Field lex.Token
}

var _ Expr = FieldExpr{}

func (expr FieldExpr) Pos() ast.Position {
	return expr.Expr.Pos().Merge(expr.Field)
}

func (expr FieldExpr) Format(w io.Writer) {
	expr.Expr.Format(w)
	_, _ = io.WriteString(w, ".")
	_, _ = io.WriteString(w, expr.Field.Value)
}

func (expr FieldExpr) Children() Children {
	return Children{expr.Expr, expr.Field}
}

func (FieldExpr) isExpr() {}
func (FieldExpr) isNode() {}

type MatchesExpr struct {
	ast.Position
	Expr      Expr
	MatchesKw lex.Token
	Pat       Pat
}

var _ Expr = MatchesExpr{}

func (expr MatchesExpr) Children() Children {
	return Children{expr.Expr, expr.MatchesKw, expr.Pat}
}

func (MatchesExpr) isExpr() {}
func (MatchesExpr) isNode() {}

func (expr MatchesExpr) Format(w io.Writer) {
	expr.Expr.Format(w)
	_, _ = io.WriteString(w, " matches ")
	expr.Pat.Format(w)
}

type ListExpr struct {
	ast.Position
	OpenBracket  lex.Token
	Exprs        []Expr
	Commas       []lex.Token // alternates with exprs, there may or may not be a trailing comma
	CloseBracket lex.Token
}

var _ Expr = ListExpr{}

func (expr ListExpr) Children() Children {
	children := Children{expr.OpenBracket}
	for i, elem := range expr.Exprs {
		children = append(children, elem)
		if i < len(expr.Commas) {
			children = append(children, expr.Commas[i])
		}
	}
	children = append(children, expr.CloseBracket)
	return children
}

func (ListExpr) isExpr() {}
func (ListExpr) isNode() {}

func (expr ListExpr) Format(w io.Writer) {
	_, _ = io.WriteString(w, "[")

	for i, expr := range expr.Exprs {
		if i > 0 {
			_, _ = io.WriteString(w, ", ")
		}
		expr.Format(w)
	}

	_, _ = io.WriteString(w, "]")
}

type ObjectExpr struct {
	ast.Position
	OpenBrace  lex.Token
	Fields     *orderedmap.OrderedMap[lex.Token, Expr]
	Commas     []lex.Token // alternates with fields, there may or may not be a trailing comma
	CloseBrace lex.Token
	Base       Expr // may be nil, the `{ ...expr }` field if one exists
}

var _ Expr = ObjectExpr{}

func (expr ObjectExpr) Children() Children {
	children := make(Children, 0, expr.Fields.Len()*3+2)
	var i int
	children = append(children, expr.OpenBrace)
	for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
		if entry.Key.Position != entry.Value.Pos() {
			// avoid doubling up due to object punning
			children = append(children, entry.Key)
		}
		children = append(children, entry.Value)
		if i < len(expr.Commas) {
			children = append(children, expr.Commas[i])
			i++
		}
	}

	if expr.Base != nil {
		children = append(children, expr.Base)
	}

	return append(children, expr.CloseBrace)
}

func (ObjectExpr) isExpr() {}
func (ObjectExpr) isNode() {}

func (expr ObjectExpr) Format(w io.Writer) {
	if expr.Fields.Len() == 0 {
		_, _ = io.WriteString(w, "{}")
		return
	}
	_, _ = io.WriteString(w, "{")

	for entry := expr.Fields.Oldest(); entry != nil; entry = entry.Next() {
		_, _ = io.WriteString(w, " ")
		_, _ = io.WriteString(w, entry.Key.Value)
		_, _ = io.WriteString(w, ": ")
		entry.Value.Format(w)
	}

	_, _ = io.WriteString(w, " }")
}

type BinaryExpr struct {
	Op    lex.Token
	Left  Expr
	Right Expr
}

var _ Expr = BinaryExpr{}

func (expr BinaryExpr) Pos() ast.Position {
	return expr.Left.Pos().Merge(expr.Right)
}

func (expr BinaryExpr) Children() Children {
	return Children{expr.Left, expr.Op, expr.Right}
}

func (BinaryExpr) isExpr() {}
func (BinaryExpr) isNode() {}

func (expr BinaryExpr) Format(w io.Writer) {
	expr.Left.Format(w)
	_, _ = io.WriteString(w, " ")
	_, _ = io.WriteString(w, expr.Op.String())
	_, _ = io.WriteString(w, " ")
	expr.Right.Format(w)
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

func (expr UnaryExpr) Format(w io.Writer) {
	_, _ = io.WriteString(w, expr.Op.String())
	expr.Expr.Format(w)
}

type TryExpr struct {
	TryKw lex.Token
	Expr  Expr
}

var _ Expr = TryExpr{}

func (expr TryExpr) Children() Children {
	return Children{expr.TryKw, expr.Expr}
}

func (expr TryExpr) Pos() ast.Position {
	return expr.TryKw.Pos().Merge(expr.Expr)
}

func (TryExpr) isExpr() {}
func (TryExpr) isNode() {}

func (expr TryExpr) Format(w io.Writer) {
	_, _ = io.WriteString(w, "try ")
	expr.Expr.Format(w)
}
