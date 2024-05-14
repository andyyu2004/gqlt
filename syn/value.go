package syn

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/gqlparser/lexer"
	"github.com/movio/gqlt/internal/lex"
)

type ValueKind int

const (
	Variable ValueKind = iota
	IntValue
	FloatValue
	StringValue
	BlockValue
	BooleanValue
	NullValue
	EnumValue
	ListValue
	ObjectValue
)

type Value struct {
	Raw      string
	Fields   ChildValueList
	Kind     ValueKind
	Position ast.Position `dump:"-"`
	Comment  *CommentGroup

	// Require validation
	Definition         *Definition
	VariableDefinition *VariableDefinition
	ExpectedType       *Type
}

// Children implements Node.
func (v *Value) Children() Children {
	if v == nil {
		return nil
	}

	switch v.Kind {
	case ListValue, ObjectValue:
		children := Children{}
		for _, child := range v.Fields {
			children = append(children, child)
		}
		return children
	default:
		if v.Fields != nil {
			panic(fmt.Errorf("unexpected fields for value kind %d", v.Kind))
		}

		var kind lex.TokenKind
		switch v.Kind {
		case Variable:
			kind = lex.Name
		case IntValue:
			kind = lex.Int
		case FloatValue:
			kind = lex.Float
		case StringValue:
			kind = lex.String
		case BlockValue:
			kind = lex.BlockString
		case NullValue:
			kind = lex.Null
		case EnumValue:
			// be nice to give the enum value a special type maybe for highlighting
			kind = lex.Name
		case BooleanValue:
			switch v.Raw {
			case "true":
				kind = lex.True
			case "false":
				kind = lex.False
			default:
				panic(fmt.Errorf("unexpected boolean value %q", v.Raw))
			}
		default:
			panic(fmt.Errorf("unexpected value kind %d", v.Kind))
		}

		return Children{
			lex.Token{
				Kind:     kind,
				Value:    v.Raw,
				Position: v.Position,
			},
		}

	}
}

func (v *Value) Format(io.Writer) {}

func (v *Value) Pos() ast.Position {
	return v.Position
}

func (*Value) isNode() {}

var _ Node = new(Value)

type ChildValue struct {
	Name     lexer.Token
	Value    *Value
	Position ast.Position `dump:"-"`
	Comment  *CommentGroup
}

func (v *ChildValue) Children() Children {
	return Children{v.Name, v.Value}
}

func (*ChildValue) Format(io.Writer) {}

func (v *ChildValue) Pos() ast.Position {
	return v.Position
}

func (*ChildValue) isNode() {}

var _ Node = new(ChildValue)

func (v *Value) Value(vars map[string]interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch v.Kind {
	case Variable:
		if value, ok := vars[v.Raw]; ok {
			return value, nil
		}
		if v.VariableDefinition != nil && v.VariableDefinition.DefaultValue != nil {
			return v.VariableDefinition.DefaultValue.Value(vars)
		}
		return nil, nil
	case IntValue:
		return strconv.ParseInt(v.Raw, 10, 64)
	case FloatValue:
		return strconv.ParseFloat(v.Raw, 64)
	case StringValue, BlockValue, EnumValue:
		return v.Raw, nil
	case BooleanValue:
		return strconv.ParseBool(v.Raw)
	case NullValue:
		return nil, nil
	case ListValue:
		var val []interface{}
		for _, elem := range v.Fields {
			elemVal, err := elem.Value.Value(vars)
			if err != nil {
				return val, err
			}
			val = append(val, elemVal)
		}
		return val, nil
	case ObjectValue:
		val := map[string]interface{}{}
		for _, elem := range v.Fields {
			elemVal, err := elem.Value.Value(vars)
			if err != nil {
				return val, err
			}
			val[elem.Name.Value] = elemVal
		}
		return val, nil
	default:
		panic(fmt.Errorf("unknown value kind %d", v.Kind))
	}
}

func (v *Value) String() string {
	if v == nil {
		return "<nil>"
	}
	switch v.Kind {
	case Variable:
		return "$" + v.Raw
	case IntValue, FloatValue, EnumValue, BooleanValue, NullValue:
		return v.Raw
	case StringValue, BlockValue:
		return strconv.Quote(v.Raw)
	case ListValue:
		var val []string
		for _, elem := range v.Fields {
			val = append(val, elem.Value.String())
		}
		return "[" + strings.Join(val, ",") + "]"
	case ObjectValue:
		var val []string
		for _, elem := range v.Fields {
			val = append(val, elem.Name.Value+":"+elem.Value.String())
		}
		return "{" + strings.Join(val, ",") + "}"
	default:
		panic(fmt.Errorf("unknown value kind %d", v.Kind))
	}
}

func (v *Value) Dump() string {
	return v.String()
}
