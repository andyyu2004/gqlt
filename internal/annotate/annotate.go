package annotate

import (
	"fmt"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

type Annotation interface {
	ast.HasPosition
	Message() string
}

func Annotate[S ~[]T, T Annotation](src string, annotations S) string {
	lines := strings.Split(src, "\n")
	// Process annotations in reverse order to avoid shifting line numbers
	for i := len(annotations) - 1; i >= 0; i-- {
		annotation := annotations[i]
		pos := annotation.Pos()
		len := pos.End - pos.Start
		padding := ""
		if pos.Column > 1 {
			// one space for the comment character and one due to 1-indexing
			padding = strings.Repeat(" ", pos.Column-2)
		}
		caret := fmt.Sprintf("%s%s", padding, strings.Repeat("^", len))
		annotationComment := fmt.Sprintf("#%s %s", caret, annotation.Message())
		lines = append(lines[:pos.Line], append([]string{annotationComment}, lines[pos.Line:]...)...)
	}

	return strings.Join(lines, "\n")
}
