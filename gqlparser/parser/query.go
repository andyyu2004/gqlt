package parser

import (
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/movio/gqlt/internal/lex"

	//nolint:revive
	. "github.com/movio/gqlt/syn"
)

func ParseQuery(source *ast.Source) (*QueryDocument, error) {
	p := New(lexer.New(source))
	return p.parseQueryDocument(), p.err
}

func (p *parser) parseQueryDocument() *QueryDocument {
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
				doc.Fragments = append(doc.Fragments, p.ParseFragmentDefinition())
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

func (p *parser) ParseOperationDefinition() *OperationDefinition {
	if p.peek().Kind == lexer.BraceL {
		return &OperationDefinition{
			Position:       p.peekPos(),
			Comment:        p.comment,
			Operation:      Query,
			OperationToken: nil,
			SelectionSet:   p.parseRequiredSelectionSet(),
		}
	}

	var od OperationDefinition
	od.Position = p.peekPos()
	od.Comment = p.comment
	od.Operation, od.OperationToken = p.parseOperationType()

	if p.peek().Kind == lexer.Name {
		od.Name = p.next()
	}

	od.VariableDefinitions = p.parseVariableDefinitions()
	od.Directives = p.parseDirectives(false)
	od.SelectionSet = p.parseRequiredSelectionSet()

	return &od
}

func (p *parser) parseOperationType() (Operation, *lex.Token) {
	tok := p.next()
	token := lex.Token{Value: tok.Value, Position: tok.Position}
	switch tok.Value {
	case "query":
		token.Kind = lex.Query
		return Query, &token
	case "mutation":
		token.Kind = lex.Mutation
		return Mutation, &token
	case "subscription":
		token.Kind = lex.Subscription
		return Subscription, &token
	}
	p.unexpectedToken(tok)
	return "", nil
}

func (p *parser) parseVariableDefinitions() VariableDefinitionList {
	var defs []*VariableDefinition
	p.many(lexer.ParenL, lexer.ParenR, func() {
		defs = append(defs, p.parseVariableDefinition())
	})

	return defs
}

func (p *parser) parseVariableDefinition() *VariableDefinition {
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

func (p *parser) parseVariable() lexer.Token {
	p.expect(lexer.Dollar)
	return p.parseName()
}

func (p *parser) parseOptionalSelectionSet() SelectionSet {
	var selections []Selection
	p.some(lexer.BraceL, lexer.BraceR, func() {
		selections = append(selections, p.parseSelection())
	})

	return selections
}

func (p *parser) parseRequiredSelectionSet() SelectionSet {
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

func (p *parser) parseSelection() Selection {
	if p.peek().Kind == lexer.Spread {
		return p.parseFragment()
	}
	return p.parseField()
}

func (p *parser) parseField() *Field {
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

func (p *parser) parseArguments(isConst bool) ArgumentList {
	var arguments ArgumentList
	p.many(lexer.ParenL, lexer.ParenR, func() {
		arguments = append(arguments, p.parseArgument(isConst))
	})

	return arguments
}

func (p *parser) parseArgument(isConst bool) *Argument {
	arg := Argument{}
	arg.Position = p.peekPos()
	arg.Comment = p.comment
	arg.Name = p.parseName()
	p.expect(lexer.Colon)

	arg.Value = p.parseValueLiteral(isConst)
	return &arg
}

func (p *parser) parseFragment() Selection {
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
		onKw := p.next() // "on"
		def.OnKw = lex.Token{
			Kind:     lex.On,
			Value:    onKw.Value,
			Position: onKw.Position,
		}

		def.TypeCondition = p.parseTypeName()
	}

	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *parser) ParseFragmentDefinition() *FragmentDefinition {
	var def FragmentDefinition
	def.Position = p.peekPos()
	def.Comment = p.comment
	def.FragmentKw = lex.Token{
		Kind:     lex.Fragment,
		Value:    p.peek().Value,
		Position: p.peek().Position,
	}
	p.expectKeyword("fragment")

	def.Name = p.parseFragmentName()
	def.VariableDefinition = p.parseVariableDefinitions()

	onKw, _ := p.expectKeyword("on")
	def.OnKw = lex.Token{
		Kind:     lex.On,
		Value:    onKw.Value,
		Position: onKw.Position,
	}

	def.TypeCondition = p.parseTypeName()
	def.Directives = p.parseDirectives(false)
	def.SelectionSet = p.parseRequiredSelectionSet()
	return &def
}

func (p *parser) parseFragmentName() lex.Token {
	if p.peek().Value == "on" {
		p.unexpectedError()
		return lex.Token{}
	}

	// treat a fragment name as a type name (for highlighting, maybe there's another more appropriate token type?)
	return p.parseTypeName()
}

func (p *parser) parseValueLiteral(isConst bool) *Value {
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
		return &Value{Position: token.Position, Comment: p.comment, Raw: p.parseVariable().Value, Kind: Variable}
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

	return &Value{Position: token.Position, Comment: p.comment, Raw: raw, Kind: kind}
}

func (p *parser) parseList(isConst bool) *Value {
	var values ChildValueList
	pos := p.peekPos()
	comment := p.comment
	p.many(lexer.BracketL, lexer.BracketR, func() {
		values = append(values, &ChildValue{Value: p.parseValueLiteral(isConst)})
	})

	return &Value{Fields: values, Kind: ListValue, Position: pos, Comment: comment}
}

func (p *parser) parseObject(isConst bool) *Value {
	var fields ChildValueList
	pos := p.peekPos()
	comment := p.comment
	p.many(lexer.BraceL, lexer.BraceR, func() {
		fields = append(fields, p.parseObjectField(isConst))
	})

	return &Value{Fields: fields, Kind: ObjectValue, Position: pos, Comment: comment}
}

func (p *parser) parseObjectField(isConst bool) *ChildValue {
	field := ChildValue{}
	field.Position = p.peekPos()
	field.Comment = p.comment
	field.Name = p.parseName()

	p.expect(lexer.Colon)

	field.Value = p.parseValueLiteral(isConst)
	return &field
}

func (p *parser) parseDirectives(isConst bool) []*Directive {
	var directives []*Directive

	for p.peek().Kind == lexer.At {
		if p.err != nil {
			break
		}
		directives = append(directives, p.parseDirective(isConst))
	}
	return directives
}

func (p *parser) parseDirective(isConst bool) *Directive {
	p.expect(lexer.At)

	return &Directive{
		Position:  p.peekPos(),
		Name:      p.parseName().Value,
		Arguments: p.parseArguments(isConst),
	}
}

func (p *parser) parseTypeReference() *Type {
	var typ Type

	if p.skip(lexer.BracketL) {
		typ.Position = p.peekPos()
		typ.Elem = p.parseTypeReference()
		p.expect(lexer.BracketR)
	} else {
		typ.Position = p.peekPos()
		typ.NamedType = p.parseTypeName()
	}

	if p.skip(lexer.Bang) {
		typ.NonNull = true
	}
	return &typ
}

func (p *parser) parseName() lexer.Token {
	token, _ := p.expect(lexer.Name)
	return token
}

// same as `parseName` except hints that the name refers to a type (for highlights)
func (p *parser) parseTypeName() lex.Token {
	token := p.parseName()
	return lex.Token{
		Kind:     lex.TypeName,
		Value:    token.Value,
		Position: token.Position,
	}
}
