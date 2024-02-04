package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/typecheck"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Snapshot) Hover(uri string, position protocol.Position) *protocol.Hover {
	s.log.Debugf("hover %s %v", uri, position)
	ast := s.Parse(uri).Ast
	mapper := s.Mapper(uri)
	point := protoToPoint(mapper, position)
	if point == nil {
		return nil
	}

	cursor := syn.NewCursor(ast)
	typeInfo := s.Typecheck(uri)
	if r, ty := hover(cursor, typeInfo, *point); ty != nil {
		return &protocol.Hover{
			// using typescript for syntax highlighting as we use the same syntax for types
			Contents: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "```typescript\n" + ty.String() + "\n```",
			},
			Range: lib.Ref(posToProto(mapper, r)),
		}
	}

	return nil
}

func hover(cursor *syn.Cursor[syn.Node], typeInfo typecheck.Info, point ast.Point) (syn.Node, typecheck.Ty) {
	// - Find the token that contains the given point
	// - Traverse the parent nodes to find a node type that we know how to handle (e.g. syn.Expr)
	tokenCursor := cursor.TokenAt(point)
	if tokenCursor == nil {
		return nil, nil
	}

	node := tokenCursor.Parent()
	for node != nil {
		switch node := node.Value().(type) {
		case *syn.NamePat:
			if ty, ok := typeInfo.BindingTypes[node]; ok {
				return node, ty
			}

		case syn.Expr:
			if ty, ok := typeInfo.ExprTypes[node]; ok {
				return node, ty
			}
		}

		node = node.Parent()
	}

	return nil, nil
}
