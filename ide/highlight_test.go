package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/ide"
)

func testHighlight(t *testing.T, name, content string, expect expect.Expectation) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		const path = "test.gqlt"

		changes := ide.Changes{
			ide.SetFileContent{Path: path, Content: content},
		}

		ide := ide.New()
		ide.Apply(changes)

		highlights := ide.Highlight(path)
		expect.AssertEqual(t, highlights.String())
	})
}

func TestHighlight(t *testing.T) {
	testHighlight(t, "simple", `let x = 5; let y = "test"`, expect.Expect(`1:1..1:4: keyword
1:5..1:6: variable
1:7..1:8: operator
1:9..1:10: number
1:12..1:15: keyword
1:16..1:17: variable
1:18..1:19: operator
1:21..1:27: string
`))

	testHighlight(t, "objects", `let { x } = { x: 4 }`, expect.Expect(`1:1..1:4: keyword
1:7..1:8: property
1:7..1:8: property
1:11..1:12: operator
1:15..1:16: property
1:18..1:19: number
`))
}
