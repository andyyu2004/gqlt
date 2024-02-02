package syn

import (
	"io"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
)

type FieldList []*FieldDefinition

func (l FieldList) ForName(name string) *FieldDefinition {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type EnumValueList []*EnumValueDefinition

func (l EnumValueList) ForName(name string) *EnumValueDefinition {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type DirectiveList []*Directive

func (l DirectiveList) ForName(name string) *Directive {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

func (l DirectiveList) ForNames(name string) []*Directive {
	resp := []*Directive{}
	for _, it := range l {
		if it.Name == name {
			resp = append(resp, it)
		}
	}
	return resp
}

type OperationList []*OperationDefinition

func (l OperationList) ForName(name string) *OperationDefinition {
	if name == "" && len(l) == 1 {
		return l[0]
	}
	for _, it := range l {
		if it.Name.Value == name {
			return it
		}
	}
	return nil
}

type FragmentDefinitionList []*FragmentDefinition

func (l FragmentDefinitionList) ForName(name string) *FragmentDefinition {
	for _, it := range l {
		if it.Name.Value == name {
			return it
		}
	}
	return nil
}

type VariableDefinitionList []*VariableDefinition

var _ Node = VariableDefinitionList{}

func (l VariableDefinitionList) Children() Children {
	children := make(Children, 0, len(l))
	for _, arg := range l {
		children = append(children, arg)
	}

	return children
}

func (VariableDefinitionList) Format(io.Writer) {}

func (l VariableDefinitionList) Pos() ast.Position {
	fst := l[0]
	lst := l[len(l)-1]
	return fst.Pos().Merge(lst)
}

func (VariableDefinitionList) isNode() {}

func (l VariableDefinitionList) ForName(name string) *VariableDefinition {
	for _, it := range l {
		if it.Variable.Value == name {
			return it
		}
	}
	return nil
}

type ArgumentList []*Argument

var _ Node = ArgumentList{}

func (args ArgumentList) Children() Children {
	children := make(Children, 0, len(args))
	for _, arg := range args {
		children = append(children, arg)
	}
	return children
}

func (ArgumentList) Format(io.Writer) {
}

func (l ArgumentList) Pos() ast.Position {
	fst := l[0]
	lst := l[len(l)-1]
	return fst.Pos().Merge(lst)
}

func (ArgumentList) isNode() {}

func (l ArgumentList) ForName(name string) *Argument {
	for _, it := range l {
		if it.Name.Value == name {
			return it
		}
	}
	return nil
}

type ArgumentDefinitionList []*ArgumentDefinition

func (l ArgumentDefinitionList) ForName(name string) *ArgumentDefinition {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type SchemaDefinitionList []*SchemaDefinition

type DirectiveDefinitionList []*DirectiveDefinition

func (l DirectiveDefinitionList) ForName(name string) *DirectiveDefinition {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type DefinitionList []*Definition

func (l DefinitionList) ForName(name string) *Definition {
	for _, it := range l {
		if it.Name == name {
			return it
		}
	}
	return nil
}

type OperationTypeDefinitionList []*OperationTypeDefinition

func (l OperationTypeDefinitionList) ForType(name string) *OperationTypeDefinition {
	for _, it := range l {
		if it.Type == name {
			return it
		}
	}
	return nil
}

type ChildValueList []*ChildValue

func (v ChildValueList) ForName(name string) *Value {
	for _, f := range v {
		if f.Name.Value == name {
			return f.Value
		}
	}
	return nil
}
