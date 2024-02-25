package ide_test

import (
	"testing"

	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/ide/fixture"
	"github.com/movio/gqlt/internal/slice"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestGotoDefinition(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			"name ref",
			`
let foo = 1
#   ...
1 + foo
#   ^^^
`,
		},
		{
			"object destructure",
			`
let { foo } = { foo: 1 }
#     ...
1 + foo
#   ^^^
`,
		},
		{
			"object rest destructure",
			`
let { ...foo } = { bar: 2, foo: 1 }
#        ...
1 + foo
#   ^^^
`,
		},
		{
			"list destructure",
			`
let [foo] = [1]
#    ...
1 + foo
#   ^^^
`,
		},
		{
			"nested destructure",
			`
let [a, { b, c }] = [15, { b: 16, c: 17 }];
#    .
assert [a, b, c] == [15, 16, 17]
#       ^
`,
		},

		{
			"object spread name ref",
			`
let obj = { a: 1, b: 2, c: 3 }
#   ...
let b = { ...obj }
#            ^^^
			`,
		},

		{
			"match condition",
			`
let x = 1
assert x matches foo if foo == 1
#                ...    ^^^
			`,
		},

		{
			"match condition shadowed",
			`
let foo = 1
assert foo matches foo if foo == 1
#                  ...    ^^^
			`,
		},

		{
			"match condition outer scope",
			`
let x = 1
#   .
assert foo matches _ if x == 1
#                       ^
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ide.TestWith(t, test.content, func(uri string, s ide.Snapshot) {
				fixture := fixture.Parse(test.content)
				require.Len(t, fixture.Ranges, 1)
				require.NotEmpty(t, fixture.Points)

				for _, point := range fixture.Points {
					locs := slice.Map(s.Definition(uri, point), func(loc protocol.Location) protocol.Range {
						require.Equal(t, uri, loc.URI, "definition should always be in the same file")
						return loc.Range
					})
					require.Equal(t, fixture.Ranges, locs)
				}
			})
		})
	}
}
