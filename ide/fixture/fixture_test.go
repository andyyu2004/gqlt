package fixture_test

import (
	"testing"

	"github.com/andyyu2004/gqlt/ide/fixture"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestParseFixture(t *testing.T) {
	t.Parallel()

	check := func(content string, expected fixture.Fixture) {
		require.Equal(t, expected, fixture.Parse(content))
	}

	check(`^^
test
^^ ^
		`, fixture.Fixture{
		Points: []fixture.Point{
			{Line: 1, Character: 0},
			{Line: 1, Character: 1},
			{Line: 1, Character: 3},
		},
	})

	check(`#.....
test
#.....
		`, fixture.Fixture{
		Ranges: []fixture.Range{
			{
				Start: protocol.Position{
					Line:      1,
					Character: 1,
				},
				End: protocol.Position{
					Line:      1,
					Character: 6,
				},
			},
		},
	})
}
