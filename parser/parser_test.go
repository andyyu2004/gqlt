package parser_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/andyyu2004/gqlt"
	"github.com/andyyu2004/gqlt/parser"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	const testpath = "tests"
	paths, err := gqlt.Discover(testpath)
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	for _, path := range paths {
		path := path
		idx := strings.Index(path, testpath)
		require.True(t, idx != -1)
		name := path[idx+len(testpath)+1 : len(path)-len(gqlt.Ext)]

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source, err := os.ReadFile(path)
			require.NoError(t, err)

			src := string(source)
			p, err := parser.New(&ast.Source{Name: name, Input: src})
			require.NoError(t, err)

			file, err := p.Parse()
			if err != nil {
				annotated := annotate(src, err.(parser.Errors))
				snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, annotated)
			} else {
				buf := bytes.Buffer{}
				file.Format(&buf)
				snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, buf.String())
			}
		})
	}
}

type Annotation interface {
	ast.HasPosition
	Message() string
}

func annotate[S ~[]T, T Annotation](src string, annotations S) string {
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
		annotationComment := fmt.Sprintf("#%s %s\n", caret, annotation.Message())
		lines = append(lines[:pos.Line], append([]string{annotationComment}, lines[pos.Line:]...)...)
	}

	return strings.Join(lines, "\n")
}
