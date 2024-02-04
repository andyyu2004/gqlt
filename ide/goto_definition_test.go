package ide_test

import (
	"testing"

	"github.com/andyyu2004/gqlt/ide"
	"github.com/andyyu2004/gqlt/ide/fixture"
	"github.com/andyyu2004/gqlt/internal/slice"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestGotoDefinition(t *testing.T) {
	check := func(content string) {
		ide.TestWith(t, content, func(uri string, s ide.Snapshot) {
			fixture := fixture.Parse(content)
			require.NotEmpty(t, fixture.Ranges)
			require.NotEmpty(t, fixture.Points, 1)

			for _, point := range fixture.Points {
				locs := slice.Map(s.Definition(uri, point), func(loc protocol.Location) protocol.Range {
					require.Equal(t, uri, loc.URI, "definition should always be in the same file")
					return loc.Range
				})
				require.Equal(t, fixture.Ranges, locs)
			}
		})
	}

	check(`
let foo = 1
#   ...
1 + foo
#   ^^^
`)

	check(`
let { foo } = { foo: 1 }
#     ...
1 + foo
#   ^^^
`)

	check(`
let { ...foo } = { bar: 2, foo: 1 }
#        ...
1 + foo
#   ^^^
`)

	check(`
let [foo] = [1]
#    ...
1 + foo
#   ^^^
`)

	check(`
let [a, { b, c }] = [15, { b: 16, c: 17 }];
#    .
assert [a, b, c] == [15, 16, 17]
#       ^`)
}
