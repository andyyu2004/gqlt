package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
)

// Test type errors
// If you want to check that a program compiles without errors, just put it in `tests/*.gqlt`
func TestTypecheckErrors(t *testing.T) {
	testDiagnostics(t, []diagnosticTestCase{
		{
			"set namespace",
			`
set namespace = "foo"
set namespace ["a", "b", "c"]
set namespace false`,
			expect.Expect(`
set namespace = "foo"
set namespace ["a", "b", "c"]
set namespace false
#^^^^^^^^^^^^^^^^^^^ expected string or list of strings as value for "namespace", found bool`),
		},
		{
			"multiply operator",
			`
[] * 2; # empty tuple
[1, 2, 3] * 2; # lists
[1, "foo", 3] * 2 # tuples
`,
			expect.Expect(`
[] * 2; # empty tuple
[1, 2, 3] * 2; # lists
[1, "foo", 3] * 2 # tuples
`),
		},
		{
			"add operator",
			`
[] + 2;
[1] + [false];
[1] + [1]
`,
			expect.Expect(`
[] + 2;
#^^^^^^ cannot apply '+' to '[]' and 'number'
[1] + [false];
#^^^^^^^^^^^^^ cannot append 'number[]' to 'bool[]'
[1] + [1]
`),
		},

		{
			"list patterns",
			`
let [a, b] = {}
let [a, b] = []
let [a, b] = [1]
let [a, b] = [1, 2]
let [a, b] = [1, 2]
let [a, b] = [1, 2, 3]
`,
			expect.Expect(`
let [a, b] = {}
#   ^^^^^^ cannot bind {} to a list pattern
let [a, b] = []
#   ^^^^^^ cannot bind tuple with 2 elements to a pattern with 0 elements
let [a, b] = [1]
let [a, b] = [1, 2]
let [a, b] = [1, 2]
let [a, b] = [1, 2, 3]
`),
		},
		{
			"object patterns",
			`
let {} = []
let [a, b, c] = [1, 2, 3]
let { a, b, c } = { a, b }
let { a, a, b } = { a, b }
`,
			expect.Expect(`
let {} = []
#   ^^ cannot bind [] to an object pattern
let [a, b, c] = [1, 2, 3]
let { a, b, c } = { a, b }
#   ^^^^^^^^^^^ field 'c' not found in object
let { a, a, b } = { a, b }
#   ^^^^^^^^^^^ field 'a' specified twice
`),
		},
		{
			"missing subselection", `
let x = query { foo }
		#   ^`, expect.Expect(`
let x = query { foo }
#               ^^^ field 'foo' of type '{ id: ID, any: Any, string: String, int: Int, float: Float, boolean: Boolean }' must have a selection of subfields
		#   ^`),
		},

		{
			"list missing subselection", `
let x = query { foos }
		#   ^`, expect.Expect(`
let x = query { foos }
#               ^^^^ field 'foos' of type '{ id: ID, any: Any, string: String, int: Int, float: Float, boolean: Boolean }[]' must have a selection of subfields
		#   ^`),
		},

		{
			"query missing field", `
let x = query { nonexistent }
`, expect.Expect(`
let x = query { nonexistent }
#               ^^^^^^^^^^^ field 'nonexistent' does not exist on type '{ foo: Foo, foos: { id: ID, any: Any, string: String, int: Int, float: Float, boolean: Boolean }[], fail: Int, animals: AnimalQuery, recursive: Recursive, __schema: __Schema, __type: __Type }'
`),
		},
		{
			"regex match", `
let x = "foo" =~ 2
let x = 2 =~ "foo"
let x = "foo" !~ 2
let x = "foo" !~ "foo"
let x = 2 !~ "foo"
`, expect.Expect(`
let x = "foo" =~ 2
#       ^^^^^^^^^^ cannot apply '=~' to 'string' and 'number' (expected string and string)
let x = 2 =~ "foo"
#       ^^^^^^^^^^ cannot apply '=~' to 'number' and 'string' (expected string and string)
let x = "foo" !~ 2
#       ^^^^^^^^^^ cannot apply '!~' to 'string' and 'number' (expected string and string)
let x = "foo" !~ "foo"
let x = 2 !~ "foo"
#       ^^^^^^^^^^ cannot apply '!~' to 'number' and 'string' (expected string and string)
`),
		},
		{
			"object base not an object", `
let x = 2
let y = { ...x }
			`, expect.Expect(`
let x = 2
let y = { ...x }
#       ^^^^^^^^ object base must be an object, got 'number'
			`),
		},
	}...)
}
