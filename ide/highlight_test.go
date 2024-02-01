package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/ide"
)

func TestHighlight(t *testing.T) {
	check := func(name, content string, expect expect.Expectation) {
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

	check("simple", `let x = 5; let y = "test"`, expect.Expect(`1:1..1:4: keyword
1:5..1:6: variable
1:7..1:8: operator
1:9..1:10: number
1:12..1:15: keyword
1:16..1:17: variable
1:18..1:19: operator
1:21..1:27: string
`))

	check("objects", `let { x } = { x: 4 }`, expect.Expect(`1:1..1:4: keyword
1:7..1:8: property
1:7..1:8: property
1:11..1:12: operator
1:15..1:16: property
1:18..1:19: number
`))

	check("literal patterns", `let 42 = 42`, expect.Expect(`1:1..1:4: keyword
1:5..1:7: number
1:8..1:9: operator
1:10..1:12: number
`))

	check("queries", `let x = query { foo bar }
mutation { bar }`, expect.Expect(`1:1..1:4: keyword
1:5..1:6: variable
1:7..1:8: operator
1:9..1:14: keyword
1:17..1:20: property
1:21..1:24: property
2:1..2:9: keyword
2:12..2:15: property
`))

	check("fragments", `fragment Foo on Bar { baz foo }`, expect.Expect(`1:1..1:9: keyword
1:10..1:13: type
1:14..1:16: keyword
1:17..1:20: type
1:23..1:26: property
1:27..1:30: property
`))

	check("set", `set namespace "foo/bar"`, expect.Expect(`1:1..1:4: keyword
1:5..1:14: variable
1:16..1:25: string
`))
}
