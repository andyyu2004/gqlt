package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/movio/gqlt/ide"
)

func TestHighlight(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectation expect.Expectation
	}{
		{
			"simple",
			`let x = 5; let y = "test"`,
			expect.Expect(`1:1..1:4: keyword
1:5..1:6: variable
1:7..1:8: operator
1:9..1:10: number
1:12..1:15: keyword
1:16..1:17: variable
1:18..1:19: operator
1:21..1:27: string
`),
		},
		{
			"objects",
			`let { x } = { x: 4 }`,
			expect.Expect(`1:1..1:4: keyword
1:7..1:8: property
1:11..1:12: operator
1:15..1:16: property
1:18..1:19: number
`),
		},
		{
			"literal patterns",
			`let 42 = 42`,
			expect.Expect(`1:1..1:4: keyword
1:5..1:7: number
1:8..1:9: operator
1:10..1:12: number
`),
		},
		{
			"queries",
			`let x = query { foo bar }
mutation { bar }`,
			expect.Expect(`1:1..1:4: keyword
1:5..1:6: variable
1:7..1:8: operator
1:9..1:14: keyword
1:17..1:20: property
1:21..1:24: property
2:1..2:9: keyword
2:12..2:15: property
`),
		},
		{
			"fragments",
			`fragment Foo on Bar { baz foo }`,
			expect.Expect(`1:1..1:9: keyword
1:10..1:13: type
1:14..1:16: keyword
1:17..1:20: type
1:23..1:26: property
1:27..1:30: property
`),
		},
		{
			"args",
			`query { foo(bar: "foo") }`,
			expect.Expect(`1:1..1:6: keyword
1:9..1:12: property
1:13..1:16: parameter
1:19..1:24: string
`),
		},
		{
			"fragment spread",
			`query { ...Foo }`,
			expect.Expect(`1:1..1:6: keyword
1:12..1:15: type
`),
		},
		{
			"inline fragment",
			`query { ...on Bar { bar } }`,
			expect.Expect(`1:1..1:6: keyword
1:12..1:14: keyword
1:15..1:18: type
1:21..1:24: property
`),
		},
		{
			"set",
			`set namespace "foo/bar"`,
			expect.Expect(`1:1..1:4: keyword
1:5..1:14: property
1:16..1:25: string
`),
		},
		{
			"try",
			"try query { foo }",
			expect.Expect(`1:1..1:4: keyword
1:5..1:10: keyword
1:13..1:16: property
`),
		},
		{
			"query variables",
			"query ($foo: Foo!) { foo(foo: $foo) }",
			expect.Expect(`1:1..1:6: keyword
1:9..1:12: parameter
1:14..1:17: type
1:22..1:25: property
1:26..1:29: parameter
1:31..1:32: parameter
`),
		},
		{
			"args",
			"query { foo(foo: $foo) }",
			expect.Expect(`1:1..1:6: keyword
1:9..1:12: property
1:13..1:16: parameter
1:18..1:19: parameter
`),
		},
		{
			"object patterns",
			"let { foo: bar } = { foo: 42 }",
			expect.Expect(`1:1..1:4: keyword
1:7..1:10: property
1:12..1:15: property
1:18..1:19: operator
1:22..1:25: property
1:27..1:29: number
`),
		},

		{
			"quoted object fields",
			`let { "foo": bar } = { "foo": 42 }`,
			expect.Expect(`1:1..1:4: keyword
1:8..1:13: string
1:14..1:17: property
1:20..1:21: operator
1:25..1:30: string
1:31..1:33: number
`),
		},
		{
			"matches",
			"assert 1 matches 2",
			expect.Expect(`1:1..1:7: keyword
1:8..1:9: number
1:10..1:17: keyword
1:18..1:19: number
`),
		},

		{
			"matches if",
			"assert 1 matches 2 if true",
			expect.Expect(`1:1..1:7: keyword
1:8..1:9: number
1:10..1:17: keyword
1:18..1:19: number
1:20..1:22: keyword
1:23..1:27: keyword
`),
		},

		{
			"bool arg",
			`mutation { foo(b: true) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:17: parameter
1:19..1:23: keyword
1:27..1:29: property
`),
		},

		{
			"int arg",
			`mutation { foo(b: 5) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:17: parameter
1:19..1:20: number
1:24..1:26: property
`),
		},

		{
			"string arg",
			`mutation { foo(b: 5) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:17: parameter
1:19..1:20: number
1:24..1:26: property
`),
		},

		{
			"var arg",
			`mutation { foo(x: $x) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:17: parameter
1:19..1:20: parameter
1:25..1:27: property
`),
		},

		{
			"null arg",
			`mutation { foo(x: null) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:17: parameter
1:19..1:23: keyword
1:27..1:29: property
`),
		},

		{
			"obj args",
			`mutation { foo(input: { x: 5 }) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:21: parameter
1:25..1:26: parameter
1:28..1:29: number
1:35..1:37: property
`),
		},

		{
			"obj args 2",
			`mutation { foo(input: { x: 1 y: 2 }) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:21: parameter
1:25..1:26: parameter
1:28..1:29: number
1:30..1:31: parameter
1:33..1:34: number
1:40..1:42: property
`),
		},

		{
			"obj args 3",
			`mutation { foo(input: { name: "a", ty: BOOL }) { id } }`, expect.Expect(`1:1..1:9: keyword
1:12..1:15: property
1:16..1:21: parameter
1:25..1:29: parameter
1:32..1:35: string
1:36..1:38: parameter
1:40..1:44: parameter
1:50..1:52: property
`),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			test := test
			t.Parallel()
			ide.TestWith(t, test.content, func(uri string, s ide.Snapshot) {
				test.expectation.AssertEqual(t, s.Highlight(uri).String())
			})
		})
	}
}

func TestHighlightScratch(t *testing.T) {
	s := `mutation { foo(input: { name: "a", ty: BOOL }) { id } }`
	ide.TestWith(t, s, func(uri string, s ide.Snapshot) {
		t.Log(s.Highlight(uri).String())
	})
}
