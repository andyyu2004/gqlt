package parser

import (
	"github.com/andyyu2004/gqlt/gqlparser/lexer"

	//nolint:revive
	. "github.com/andyyu2004/gqlt/gqlparser/ast"
)

func ParseQuery(source *Source) (*QueryDocument, error) {
	p := Parser{
		lexer: lexer.New(source),
	}
	return p.parseQueryDocument(), p.err
}

func (p *Parser) parseQueryDocument() *QueryDocument {
	var doc QueryDocument
	for p.peek().Kind != lexer.EOF {
		if p.err != nil {
			return &doc
		}
		doc.Position = p.peekPos()
		switch p.peek().Kind {
		case lexer.Name:
			switch p.peek().Value {
			case "query", "mutation", "subscription":
				doc.Operations = append(doc.Operations, p.ParseOperationDefinition())
			case "fragment":
				doc.Fragments = append(doc.Fragments, p.parseFragmentDefinition())
			default:
				p.unexpectedError()
			}
		case lexer.BraceL:
			doc.Operations = append(doc.Operations, p.ParseOperationDefinition())
		default:
			p.unexpectedError()
		}
	}

	return &doc
}

func (p *Parser) ParseOperationDefinition() *OperationDefinition {
	if p.peek().Kind == lexer.BraceL {
		return &OperationDefinition{
			Position:     p.peekPos(),
			Comment:      p.comment,
			Operation:    Query,
			SelectionSet: p.parseRequiredSelectionSet(),
		}
	}

	var od OperationDefinition
	od.Position = p.peekPos()
	od.Comment = p.comment
	od.Operation = p.parseOperationType()

	if p.peek().Kind == lexer.Name {
		od.Name = p.next().Value
	}

	od.VariableDefinitions = p.parseVariableDefinitions()
	od.Directives = p.parseDirectives(false)
	od.SelectionSet = p.parseRequiredSelectionSet()

	return &od
}

func (p *Parser) parseOperationType() Operation {
	tok := p.next()
	switch tok.Value {
	case "query":
		return Query
	case "mutation":
		return Mutation
	case "subscription":
		return Subscription
	}
	p.unexpectedToken(tok)
	return ""
}

func (p *Parser) parseVariableDefinitions() VariableDefinitionList {
	var defs []*VariableDefinition
	p.many(lexer.ParenL, lexer.ParenR, func() {
		defs = append(defs, p.parseVariableDefinition())
	})

	return defs
}

func (p *Parser) parseVariableDefinition() *VariableDefinition {
	var def VariableDefinition
	def.Position = p.peekPos()
	def.Comment = p.comment
	def.Variable = p.parseVariable()

	p.expect(lexer.Colon)

	def.Type = p.parseTypeReference()

	if p.skip(lexer.Equals) {
		def.DefaultValue = p.parseValueLiteral(true)
	}

	def.Directives = p.parseDirectives(false)

	return &def
}

func (p *Parser) parseVariable() string {
	p.expect(lexer.Dollar)
	return p.parseName()
}

func (p *Parser) parseOptionalSelectionSet() SelectionSet {
	var selections []Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *Parser) parseRequiredSelectionSet() SelectionSet {
	if p.peek().Kind != lexer.BraceL {
		p.error(p.peek(), "Expected %s, found %s", lexer.BraceL, p.peek().Kind.String())
		return nil
	}

	var selections []Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *Parser) parseSelection() Selection {
	if p.peek().Kind == lexer.Spread {
		return p.parseFragment()
	}
	return p.parseField()
}

func (p *Parser) parseField() *Field {
	var field Field
	field.Position = p.peekPos()
	field.Comment = p.comment
	field.Alias = p.parseName()

	if p.skip(lexer.Colon) {
		field.Name = p.parseName()
	} else {
		field.Name = field.Alias
	}

	field.Arguments = p.parseArguments(false)
	field.Directives = p.parseDirectives(false)
	if p.peek().Kind == lexer.BraceL {
		field.SelectionSet = p.parseOptionalSelectionSet()
	}

	return &field
}

func (p *Parser) parseArguments(isConst bool) ArgumentList {
	var arguments ArgumentList
	p.many(lexer.ParenL, lexer.ParenR, func() {
		arguments = append(arguments, p.parseArgument(isConst))
	})

	return arguments
}

func (p *Parser) parseArgument(isConst bool) *Argument {
	arg := Argument{}
	arg.Position = p.peekPos()
	arg.Comment = p.comment
	arg.Name = p.parseName()
	p.expect(lexer.Colon)

	arg.Value = p.parseValueLiteral(isConst)
	return &arg
}

func (p *Parser) parseFragment() Selection {
	_, comment := p.expect(lexer.Spread)

	if peek := p.peek(); peek.Kind == lexer.Name && peek.Value != "on" {
		return &FragmentSpread{
			Position:   p.peekPos(),
			Comment:    comment,
			Name:       p.parseFragmentName(),
			Directives: p.parseDirectives(false),
		}
	}

	var def InlineFragment
	def.Position = p.peekPos()
	def.Comment = comment
	if p.peek().Value == "on" {
		p.next() // "on"

		def.TypeCondition = p.parseName()
	}

	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *Parser) parseFragmentDefinition() *FragmentDefinition {
	var def FragmentDefinition
	def.Position = p.peekPos()
	def.Comment = p.comment
	p.expectKeyword("fragment")

	def.Name = p.parseFragmentName()
	def.VariableDefinition = p.parseVariableDefinitions()

	p.expectKeyword("on")

	def.TypeCondition = p.parseName()
	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *Parser) parseFragmentName() string {
	if p.peek().Value == "on" {
		p.unexpectedError()
		return ""
	}

	return p.parseName()
}

func (p *Parser) parseValueLiteral(isConst bool) *Value {
	token := p.peek()
	raw := token.Value

	var kind ValueKind
	switch token.Kind {
	case lexer.BracketL:
		return p.parseList(isConst)
	case lexer.BraceL:
		return p.parseObject(isConst)
	case lexer.Dollar:
		if isConst {
			p.unexpectedError()
			return nil
		}
		return &Value{Position: &token.Pos, Comment: p.comment, Raw: p.parseVariable(), Kind: Variable}
	case lexer.Minus:
		p.next()
		tok := p.peek()
		switch tok.Kind {
		case lexer.Int:
			kind = IntValue
			raw = "-" + tok.Value
		case lexer.Float:
			kind = FloatValue
			raw = "-" + tok.Value
		default:
			p.unexpectedError()
		}
	case lexer.Int:
		kind = IntValue
	case lexer.Float:
		kind = FloatValue
	case lexer.String:
		kind = StringValue
	case lexer.BlockString:
		kind = BlockValue
	case lexer.Name:
		switch token.Value {
		case "true", "false":
			kind = BooleanValue
		case "null":
			kind = NullValue
		default:
			kind = EnumValue
		}
	default:
		p.unexpectedError()
		return nil
	}

	p.next()

	return &Value{Position: &token.Pos, Comment: p.comment, Raw: raw, Kind: kind}
}

func (p *Parser) parseList(isConst bool) *Value {
	var values ChildValueList
	pos := p.peekPos()
	comment := p.comment
	p.many(lexer.BracketL, lexer.BracketR, func() {
		values = append(values, &ChildValue{Value: p.parseValueLiteral(isConst)})
	})

	return &Value{Children: values, Kind: ListValue, Position: pos, Comment: comment}
}

func (p *Parser) parseObject(isConst bool) *Value {
	var fields ChildValueList
	pos := p.peekPos()
	comment := p.comment
	p.many(lexer.BraceL, lexer.BraceR, func() {
		fields = append(fields, p.parseObjectField(isConst))
	})

	return &Value{Children: fields, Kind: ObjectValue, Position: pos, Comment: comment}
}

func (p *Parser) parseObjectField(isConst bool) *ChildValue {
	field := ChildValue{}
	field.Position = p.peekPos()
	field.Comment = p.comment
	field.Name = p.parseName()

	p.expect(lexer.Colon)

	field.Value = p.parseValueLiteral(isConst)
	return &field
}

func (p *Parser) parseDirectives(isConst bool) []*Directive {
	var directives []*Directive

	for p.peek().Kind == lexer.At {
		if p.err != nil {
			break
		}
		directives = append(directives, p.parseDirective(isConst))
	}
	return directives
}

func (p *Parser) parseDirective(isConst bool) *Directive {
	p.expect(lexer.At)

	return &Directive{
		Position:  p.peekPos(),
		Name:      p.parseName(),
		Arguments: p.parseArguments(isConst),
	}
}

func (p *Parser) parseTypeReference() *Type {
	var typ Type

	if p.skip(lexer.BracketL) {
		typ.Position = p.peekPos()
		typ.Elem = p.parseTypeReference()
		p.expect(lexer.BracketR)
	} else {
		typ.Position = p.peekPos()
		typ.NamedType = p.parseName()
	}

	if p.skip(lexer.Bang) {
		typ.NonNull = true
	}
	return &typ
}

func (p *Parser) parseName() string {
	token, _ := p.expect(lexer.Name)

	return token.Value
}
