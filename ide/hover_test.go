package ide_test

import (
	"strings"
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/ide"
	"github.com/andyyu2004/gqlt/ide/fixture"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestHover(t *testing.T) {
	check := func(name, content string, expectation expect.Expectation) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ide.TestWith(t, content, func(path string, s ide.Snapshot) {
				fixture := fixture.Parse(content)
				require.Empty(t, fixture.Ranges)
				require.NotEmpty(t, fixture.Points)

				var expected string
				// expect all points to yield the same hover
				for _, point := range fixture.Points {
					hover := s.Hover(path, point)
					var result string
					if hover != nil {
						content := hover.Contents.(protocol.MarkupContent).Value
						result = strings.TrimSuffix(strings.TrimPrefix(content, "```typescript\n"), "\n```")
					} else {
						result = "no hover"
					}

					if expected == "" {
						expected = result
					}

					require.Equal(t, expected, result, "expected all points to yield the same hover")
				}

				expectation.AssertEqual(t, expected)
			})
		})
	}

	check("hover empty space", `
let x = 5
#         ^^`, expect.Expect(`no hover`))

	check("number literal type", `
let x = 5
#       ^`, expect.Expect(`number`))

	check("string literal type", `
let x = "test"
#       ^^^^^^`, expect.Expect(`string`))

	check("bool literal type", `
let x = false
#       ^^^^^`, expect.Expect(`bool`))

	check("tuple literal type", `
let x = [5, "test", false]
        ^`, expect.Expect(`[number, string, bool]`))

	check("list literal type", `
let x = [1, 2, 3]
        ^`, expect.Expect(`number[]`))

	check("object literal type", `
let x = { a: 1, b: "test", c: false }
        ^`, expect.Expect(`{ a: number, b: string, c: bool }`))

	check("name ref", `
let x = 2
1 + x
    ^`, expect.Expect(`number`))

	check("name def", `
let x = 2
    ^`, expect.Expect(`number`))
}
