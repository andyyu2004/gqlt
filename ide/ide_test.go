package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/movio/gqlt/ide"
)

func TestIDE(t *testing.T) {
	ide.TestWith(t, "let x = 5", func(uri string, s ide.Snapshot) {
		ast := s.Parse(uri).Ast
		expect.Expect(`let x = 5;
`).AssertEqual(t, ast.String())
	})
}
