package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/ide"
)

func TestIDE(t *testing.T) {
	const path = "test.gqlt"

	changes := ide.Changes{
		ide.SetFileContent{
			Path:    path,
			Content: "let x = 5",
		},
	}

	ide := ide.New()
	ide.Apply(changes)

	ast := ide.Parse(path)
	expect.Expect(`let x = 5;
`).AssertEqual(t, ast.String())
}
