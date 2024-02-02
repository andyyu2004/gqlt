package parser_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/andyyu2004/gqlt"
	"github.com/andyyu2004/gqlt/internal/annotate"
	"github.com/andyyu2004/gqlt/internal/parser"

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
			p := parser.New(&ast.Source{Name: name, Input: src})
			require.NoError(t, err)

			file, err := p.Parse()
			if err != nil {
				annotated := annotate.Annotate(src, err.(ast.Errors))
				snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, annotated)
			} else {
				buf := bytes.Buffer{}
				file.Format(&buf)
				snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, buf.String())
			}
		})
	}
}
