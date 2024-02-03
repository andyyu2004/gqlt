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

// A caret each represents a point
// A series of .... represents a range
func Parse(content string) Fixture {
	lines := strings.Split(content, "\n")
	var points []protocol.Position

	// markers on the first line are ignored as there's nothing for them to point at
	for l, line := range lines[1:] {
		for c, r := range line {
			if r == '^' {
				points = append(points, protocol.Position{Line: uint32(l), Character: uint32(c)})
			}
		}
	}

	var ranges []protocol.Range
	for l, line := range lines[1:] {
		var i int
		var isCommentLine bool
		for i < len(line) {
			if line[i] == '#' {
				isCommentLine = true
			}

			if isCommentLine && line[i] == '.' {
				start := protocol.Position{Line: uint32(l), Character: uint32(i)}
				end := protocol.Position{Line: uint32(l), Character: uint32(i + 1)}
				for i+1 < len(line) && line[i+1] == '.' {
					i++
					end.Character++
				}
				ranges = append(ranges, protocol.Range{Start: start, End: end})
			}
			i++
		}
	}

	return Fixture{Points: points, Ranges: ranges}
}
