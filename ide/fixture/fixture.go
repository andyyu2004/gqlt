package fixture

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type (
	Point = protocol.Position
	Range = protocol.Range
)

type Fixture struct {
	Points []Point
	Ranges []Range
}

// A caret represents a point
func Parse(content string) Fixture {
	lines := strings.Split(content, "\n")
	var points []protocol.Position
	// carets on the first line are ignored as they don't make sense
	for l, line := range lines[1:] {
		for c, r := range line {
			if r == '^' {
				points = append(points, protocol.Position{Line: uint32(l), Character: uint32(c)})
			}
		}
	}
	return Fixture{Points: points}
}
