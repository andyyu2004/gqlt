package parser_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"andyyu2004/gqlt/parser"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

// builtin glob doesn't implement ** :/
func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func TestParser(t *testing.T) {
	paths, err := glob("test-data", ".gqlt")
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	for _, path := range paths {
		path := path
		_, filename := filepath.Split(path)
		name := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source, err := os.ReadFile(path)
			require.NoError(t, err)

			src := string(source)
			parser, err := parser.New(&ast.Source{Name: name, Input: src})
			require.NoError(t, err)

			file, err := parser.Parse()
			require.NoError(t, err)

			buf := bytes.Buffer{}
			file.Dump(&buf)

			snaps.WithConfig(snaps.Filename(name)).MatchSnapshot(t, buf.String())
		})
	}
}
