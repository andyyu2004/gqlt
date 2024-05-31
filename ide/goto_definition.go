package ide

import (
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/typecheck"
	"github.com/movio/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Snapshot) Definition(uri string, position protocol.Position) []protocol.Location {
	s.log.Debugf("definition %s %v", uri, position)

	mapper := s.Mapper(uri)
	point := protoToPoint(mapper, position)
	if point == nil {
		return nil
	}

	typeInfo := s.Typecheck(uri)
	if node := definition(typeInfo, *point); node != nil {
		return []protocol.Location{
			{
				URI:   uri, // definitions are always in the same file currently
				Range: posToProto(mapper, node),
			},
		}
	}
	return nil
}

func definition(typeInfo typecheck.Info, point ast.Point) syn.Node {
	// - Find the token that contains the given point
	// - Traverse the parent nodes to find a `*syn.NameExpr`, `*syn.VarPat` (the only syntax that can define a name)
	// and lookup what it resolves to

	cursor := syn.NewCursor(typeInfo.Ast)
	tokenCursor := cursor.TokenAt(point)
	if tokenCursor == nil {
		return nil
	}

	node := tokenCursor.Parent()
	for node != nil {
		switch node := node.Value().(type) {
		case *syn.NamePat:
			// if on a name definition itself, return it
			return node
		case *syn.FragmentDefinition:
			// if on a fragment definition, return it
			return node
		case *syn.VarPat:
			if pat, ok := typeInfo.VarPatResolutions[node]; ok {
				return pat
			}
		case *syn.NameExpr:
			if pat, ok := typeInfo.NameResolutions[node]; ok {
				return pat
			}
		case *syn.FragmentSpread:
			if def, ok := typeInfo.FragmentResolutions[node]; ok {
				return def
			}
		}

		node = node.Parent()
	}

	return nil
}
