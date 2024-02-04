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

	check("list pattern", `
let [x] = [1, 2, 3]
     ^`, expect.Expect(`number`))

	check("list rest pattern", `
let [x, ...xs] = [1, 2, 3]
#          ^`, expect.Expect(`number[]`))

	check("tuple pattern", `
let [x, y] = [1, "s"]
#       ^`, expect.Expect(`string`))

	check("tuple rest pattern", `
let [x, ...xs] = [1, "s", false]
#          ^`, expect.Expect(`[string, bool]`))

	check("object pattern (punned)", `
let { a, b } = { a: 1, b: "s" }
#        ^`, expect.Expect(`string`))

	check("object pattern", `
let { a, b: c } = { a: 1, b: "s" }
#           ^`, expect.Expect(`string`))

	check("object rest pattern", `
let { a, ...r } = { a: 1, b: "s", c: false }
#           ^`, expect.Expect(`{ b: string, c: bool }`))

	check("matches", `
assert 1 matches 1
#        ^^^^^^^`, expect.Expect(`bool`))

	check("list index", `
let x = [1, 2, 3][0]
#   ^
`, expect.Expect(`number`))

	check("tuple literal index", `
let x = [1, "b"][1]
#   ^
`, expect.Expect(`string`))

	check("tuple variable index", `
let i = 1
let x = [1, "b"][i]
#   ^
`, expect.Expect(`any`))

	check("object literal index", `
let x = { a: 1, b: "s" }["a"]
#   ^
`, expect.Expect(`number`))

	check("object variable index", `
let k = "a"
let x = { a: 1, b: "s" }[k]
#   ^
`, expect.Expect(`any`))

	check("object field access", `
let x = { a: 1, b: "s" }.a
#   ^
`, expect.Expect(`number`))

	check("bang operator", `
let x = !1
#   ^`, expect.Expect(`bool`))

	check("not operator", `
let x = not 1 matches 1
#   ^`, expect.Expect(`bool`))

	check("negate operator", `
let y = 2
let x = -y
#   ^`, expect.Expect(`number`))

	check("minus numbers", `
let x = 1 - 1
#   ^`, expect.Expect(`number`))

	check("plus numbers", `
let x = 1 + 1
#   ^`, expect.Expect(`number`))

	check("plus strings", `
let x = "foo" + "bar"
#   ^`, expect.Expect(`string`))

	check("plus lists", `
let x = [1] + [2]
#   ^`, expect.Expect(`number[]`))

	check("multiply numbers", `
let x = 2 * 3
#   ^`, expect.Expect(`number`))

	check("multiply lists", `
let x = [1] * 3
#   ^`, expect.Expect(`number[]`))

	check("multiply tuples", `
let x = [1, "foo", false] * 2
#   ^`, expect.Expect(`[number, string, bool, number, string, bool]`))

	check("divide numbers", `
let x = 2 / 3
#   ^`, expect.Expect(`number`))
}
