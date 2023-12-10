package syn

import (
	"bytes"
	"io"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

func Dump(node Node) string {
	var buf bytes.Buffer
	node.Dump(&buf)
	return buf.String()
}

type Node interface {
	Child
	isNode()

	Dump(io.Writer)
	Children() Children
}

type File struct {
	Stmts []Stmt
}

type Children []Child

// Node | Token
type Child interface {
	ast.HasPosition
}

func (f File) Dump(w io.Writer) {
	for _, stmt := range f.Stmts {
		stmt.Dump(w)
		io.WriteString(w, ";\n")
	}
}

func (f File) String() string {
	b := strings.Builder{}
	f.Dump(&b)
	return b.String()
}
