package ide

import (
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/syn"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Snapshot) Hover(path string, position protocol.Position) *protocol.Hover {
	ast := s.Parse(path).Ast
	cursor := syn.NewCursor(ast)
	mapper := s.Mapper(path)
	point := protoToPoint(mapper, position)
	if point == nil {
		return nil
	}

	typeInfo := s.Typecheck(path)

	// - Find the token that contains the given point
	// - If the token is something we can handle, do so
	// - Otherwise, traverse the parent nodes to find a node type that we know how to handle (e.g. syn.Expr)

	tokenCursor := cursor.TokenAt(*point)
	if tokenCursor == nil {
		return nil
	}

	switch tokenCursor.Value().Kind {
	}

	par := tokenCursor.Parent()
	for par != nil {
		switch node := par.Value().(type) {
		case syn.Expr:
			if ty, ok := typeInfo.ExprTypes[node]; ok {
				return &protocol.Hover{
					// using typescript for syntax highlighting as we use the same syntax for types
					Contents: protocol.MarkupContent{
						Kind:  protocol.MarkupKindMarkdown,
						Value: "```typescript\n" + ty.String() + "\n```",
					},
					Range: lib.Ref(posToProto(mapper, node)),
				}
			}
		}

		par = par.Parent()
	}

	return nil
}
