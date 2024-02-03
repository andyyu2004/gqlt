package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/typecheck"
	"github.com/andyyu2004/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Snapshot) Definition(path string, position protocol.Position) []protocol.Location {
	s.log.Debugf("definition %s %v", path, position)

	ast := s.Parse(path).Ast
	mapper := s.Mapper(path)
	point := protoToPoint(mapper, position)
	if point == nil {
		return nil
	}

	cursor := syn.NewCursor(ast)
	typeInfo := s.Typecheck(path)
	if node := definition(cursor, typeInfo, *point); node != nil {
		return []protocol.Location{
			{
				URI:   path, // definitions are always in the same file
				Range: posToProto(mapper, node.Pos()),
			},
		}
	}
	return nil
}

func definition(cursor *syn.Cursor[syn.Node], typeInfo typecheck.Info, point ast.Point) syn.Node {
	// - Find the token that contains the given point
	// - Traverse the parent nodes to find a `*syn.NameExpr`, and lookup what it resolves to
	tokenCursor := cursor.TokenAt(point)
	if tokenCursor == nil {
		return nil
	}

	node := tokenCursor.Parent()
	for node != nil {
		switch node := node.Value().(type) {
		case *syn.NameExpr:
			if pat, ok := typeInfo.Resolutions[node]; ok {
				return pat
			}
		}

		node = node.Parent()
	}

	return nil
}
