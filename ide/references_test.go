package ide_test

import (
	"testing"

	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/ide/fixture"
	"github.com/movio/gqlt/internal/slice"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestReferences(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			"it works",
			`
let foo = 1
#   ...
print(foo)
#     ...
			`,
		},

		{
			"variable patterns",
			`
let foo = 1
#   ...
assert 1 matches $foo
#                ....
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ide.TestWith(t, test.content, func(uri string, s ide.Snapshot) {
				fixture := fixture.Parse(test.content)
				require.NotEmpty(t, fixture.Ranges)
				require.Empty(t, fixture.Points)

				for _, r := range fixture.Ranges {
					locs := slice.Map(s.References(uri, r.Start), func(loc protocol.Location) protocol.Range {
						require.Equal(t, uri, loc.URI, "reference should always be in the same file")
						return loc.Range
					})
					require.ElementsMatch(t, fixture.Ranges, locs)
				}
			})
		})
	}
}
