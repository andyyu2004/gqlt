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
				file.Dump(&buf)
				snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, buf.String())
			}
		})
	}
}

func annotate(src string, errors parser.Errors) string {
	lines := strings.Split(src, "\n")
	// Process annotations in reverse order to avoid shifting line numbers
	for i := len(errors) - 1; i >= 0; i-- {
		err := errors[i]
		len := err.End - err.Start
		padding := ""
		if err.Column > 1 {
			// one space for the comment character and one due to 1-indexing
			padding = strings.Repeat(" ", err.Column-2)
		}
		caret := fmt.Sprintf("%s%s", padding, strings.Repeat("^", len))
		annotationComment := fmt.Sprintf("#%s %s\n", caret, err.Message)
		lines = append(lines[:err.Line], append([]string{annotationComment}, lines[err.Line:]...)...)
	}

	return strings.Join(lines, "\n")
}

func insert[T any](a []T, index int, value T) []T {
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}
