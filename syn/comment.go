package syn

import (
	"strconv"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

type Comment struct {
	Value    string
	Position *ast.Position
}

func (c *Comment) Text() string {
	return strings.TrimPrefix(c.Value, "#")
}

type CommentGroup struct {
	List []*Comment
}

func (c *CommentGroup) Dump() string {
	if len(c.List) == 0 {
		return ""
	}
	var builder strings.Builder
	for _, comment := range c.List {
		builder.WriteString(comment.Value)
		builder.WriteString("\n")
	}
	return strconv.Quote(builder.String())
}
