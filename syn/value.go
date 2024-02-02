package syn

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/gqlparser/lexer"
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
	children := Children{}
	for _, f := range v.Fields {
		children = append(children, f)
	}
	return children
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
