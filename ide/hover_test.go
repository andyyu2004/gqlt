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
	tests := []struct {
		name        string
		content     string
		expectation expect.Expectation
	}{
		{
			"hover empty space",
			`
let x = 5
#         ^^`,
			expect.Expect(`no hover`),
		},
		{
			"number literal type",
			`
let x = 5
#       ^`,
			expect.Expect(`number`),
		},
		{
			"string literal type",
			`
let x = "test"
#       ^^^^^^`,
			expect.Expect(`string`),
		},
		{
			"bool literal type",
			`
let x = false
#       ^^^^^`,
			expect.Expect(`bool`),
		},
		{
			"tuple literal type",
			`
let x = [5, "test", false]
        ^`,
			expect.Expect(`[number, string, bool]`),
		},
		{
			"list literal type",
			`
let x = [1, 2, 3]
        ^`,
			expect.Expect(`number[]`),
		},
		{
			"object literal type",
			`
let x = { a: 1, b: "test", c: false }
        ^`,
			expect.Expect(`{ a: number, b: string, c: bool }`),
		},
		{
			"name ref",
			`
let x = 2
1 + x
    ^`,
			expect.Expect(`number`),
		},
		{
			"name def",
			`
let x = 2
    ^`,
			expect.Expect(`number`),
		},
		{
			"list pattern",
			`
let [x] = [1, 2, 3]
     ^`,
			expect.Expect(`number`),
		},
		{
			"list rest pattern",
			`
let [x, ...xs] = [1, 2, 3]
#          ^`,
			expect.Expect(`number[]`),
		},
		{
			"tuple pattern",
			`
let [x, y] = [1, "s"]
#       ^`,
			expect.Expect(`string`),
		},
		{
			"tuple rest pattern",
			`
let [x, ...xs] = [1, "s", false]
#          ^`,
			expect.Expect(`[string, bool]`),
		},
		{
			"object pattern (punned)",
			`
let { a, b } = { a: 1, b: "s" }
#        ^`,
			expect.Expect(`string`),
		},
		{
			"object pattern",
			`
let { a, b: c } = { a: 1, b: "s" }
#           ^`,
			expect.Expect(`string`),
		},
		{
			"object rest pattern",
			`
let { a, ...r } = { a: 1, b: "s", c: false }
#           ^`,
			expect.Expect(`{ b: string, c: bool }`),
		},
		{
			"matches",
			`
assert 1 matches 1
#        ^^^^^^^`,
			expect.Expect(`bool`),
		},
		{
			"list index",
			`
let x = [1, 2, 3][0]
#   ^
`,
			expect.Expect(`number`),
		},
		{
			"tuple literal index",
			`
let x = [1, "b"][1]
#   ^
`,
			expect.Expect(`string`),
		},
		{
			"tuple variable index",
			`
let i = 1
let x = [1, "b"][i]
#   ^
`,
			expect.Expect(`any`),
		},
		{
			"object literal index",
			`
let x = { a: 1, b: "s" }["a"]
#   ^
`,
			expect.Expect(`number`),
		},
		{
			"object variable index",
			`
let k = "a"
let x = { a: 1, b: "s" }[k]
#   ^
`,
			expect.Expect(`any`),
		},
		{
			"object field access",
			`
let x = { a: 1, b: "s" }.a
#   ^
`,
			expect.Expect(`number`),
		},
		{
			"bang operator",
			`
let x = !1
#   ^`,
			expect.Expect(`bool`),
		},
		{
			"not operator",
			`
let x = not 1 matches 1
#   ^`,
			expect.Expect(`bool`),
		},
		{
			"negate operator",
			`
let y = 2
let x = -y
#   ^`,
			expect.Expect(`number`),
		},
		{
			"minus numbers",
			`
let x = 1 - 1
#   ^`,
			expect.Expect(`number`),
		},
		{
			"plus numbers",
			`
let x = 1 + 1
#   ^`,
			expect.Expect(`number`),
		},
		{
			"plus strings",
			`
let x = "foo" + "bar"
#   ^`,
			expect.Expect(`string`),
		},
		{
			"plus lists",
			`
let x = [1] + [2]
#   ^`,
			expect.Expect(`number[]`),
		},
		{
			"multiply numbers",
			`
let x = 2 * 3
#   ^`,
			expect.Expect(`number`),
		},
		{
			"multiply lists",
			`
let x = [1] * 3
#   ^`,
			expect.Expect(`number[]`),
		},
		{
			"multiply tuples",
			`
let x = [1, "foo", false] * 2
#   ^`,
			expect.Expect(`[number, string, bool, number, string, bool]`),
		},
		{
			"divide numbers",
			`
let x = 2 / 3
#   ^`,
			expect.Expect(`number`),
		},
		{
			"query scalar", `
let x = query { int }
#   ^}`, expect.Expect(`number`),
		},
		{
			"query object field", `
let x = query { foo { id string float int boolean } }
#   ^`, expect.Expect(`{ id: string, string: string, float: number, int: number, boolean: bool }`),
		},
		{
			"query object field partial selection", `
		let x = query { foo { id boolean } }
		#   ^`, expect.Expect(`{ id: string, boolean: bool }`),
		},
		{
			"query recursive object", `
let x = query { recursive { id next { id next { id } } } }
    ^`, expect.Expect(`{ id: string, next: { id: string, next: { id: string } } }`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ide.TestWith(t, test.content, func(uri string, s ide.Snapshot) {
				fixture := fixture.Parse(test.content)
				require.Empty(t, fixture.Ranges)
				require.NotEmpty(t, fixture.Points)

				var expected string
				// expect all points to yield the same hover
				for _, point := range fixture.Points {
					hover := s.Hover(uri, point)
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

				test.expectation.AssertEqual(t, expected)
			})
		})
	}
}
