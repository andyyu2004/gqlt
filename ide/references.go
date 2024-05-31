package ide

import (
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/internal/slice"
	"github.com/movio/gqlt/internal/typecheck"
	"github.com/movio/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Snapshot) References(uri string, position protocol.Position) []protocol.Location {
	s.log.Debugf("references %s %v", uri, position)

	mapper := s.Mapper(uri)
	point := protoToPoint(mapper, position)
	if point == nil {
		return nil
	}

	typeInfo := s.Typecheck(uri)
	return slice.Map(references(typeInfo, *point), func(node syn.Node) protocol.Location {
		return protocol.Location{
			URI:   uri,
			Range: posToProto(mapper, node.Pos()),
		}
	})
}

func references(typeInfo typecheck.Info, point ast.Point) []syn.Node {
	// - Find the definition of the name at the given point
	// - Find all possible references:
	//     - If the reference resolves to the same definition, add it to the list

	def := definition(typeInfo, point)
	if def == nil {
		return nil
	}

	nodes := []syn.Node{def}

	switch def.(type) {
	case *syn.NamePat:
		for node, ref := range typeInfo.NameResolutions {
			if ref == def {
				nodes = append(nodes, node)
			}
		}

		for node, ref := range typeInfo.VarPatResolutions {
			if ref == def {
				nodes = append(nodes, node)
			}
		}
	case *syn.FragmentDefinition:
		panic("todo")
	default:
	}

	return nodes
}
