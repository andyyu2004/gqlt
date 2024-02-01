package validator

import (
	"context"
	"fmt"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/syn"
)

type Events struct {
	operationVisitor []func(walker *Walker, operation *syn.OperationDefinition)
	field            []func(walker *Walker, field *syn.Field)
	fragment         []func(walker *Walker, fragment *syn.FragmentDefinition)
	inlineFragment   []func(walker *Walker, inlineFragment *syn.InlineFragment)
	fragmentSpread   []func(walker *Walker, fragmentSpread *syn.FragmentSpread)
	directive        []func(walker *Walker, directive *syn.Directive)
	directiveList    []func(walker *Walker, directives []*syn.Directive)
	value            []func(walker *Walker, value *syn.Value)
	variable         []func(walker *Walker, variable *syn.VariableDefinition)
}

func (o *Events) OnOperation(f func(walker *Walker, operation *syn.OperationDefinition)) {
	o.operationVisitor = append(o.operationVisitor, f)
}

func (o *Events) OnField(f func(walker *Walker, field *syn.Field)) {
	o.field = append(o.field, f)
}

func (o *Events) OnFragment(f func(walker *Walker, fragment *syn.FragmentDefinition)) {
	o.fragment = append(o.fragment, f)
}

func (o *Events) OnInlineFragment(f func(walker *Walker, inlineFragment *syn.InlineFragment)) {
	o.inlineFragment = append(o.inlineFragment, f)
}

func (o *Events) OnFragmentSpread(f func(walker *Walker, fragmentSpread *syn.FragmentSpread)) {
	o.fragmentSpread = append(o.fragmentSpread, f)
}

func (o *Events) OnDirective(f func(walker *Walker, directive *syn.Directive)) {
	o.directive = append(o.directive, f)
}

func (o *Events) OnDirectiveList(f func(walker *Walker, directives []*syn.Directive)) {
	o.directiveList = append(o.directiveList, f)
}

func (o *Events) OnValue(f func(walker *Walker, value *syn.Value)) {
	o.value = append(o.value, f)
}

func (o *Events) OnVariable(f func(walker *Walker, variable *syn.VariableDefinition)) {
	o.variable = append(o.variable, f)
}

func Walk(schema *syn.Schema, document *syn.QueryDocument, observers *Events) {
	w := Walker{
		Observers: observers,
		Schema:    schema,
		Document:  document,
	}

	w.walk()
}

type Walker struct {
	Context   context.Context
	Observers *Events
	Schema    *syn.Schema
	Document  *syn.QueryDocument

	validatedFragmentSpreads map[string]bool
	CurrentOperation         *syn.OperationDefinition
}

func (w *Walker) walk() {
	for _, child := range w.Document.Operations {
		w.validatedFragmentSpreads = make(map[string]bool)
		w.walkOperation(child)
	}
	for _, child := range w.Document.Fragments {
		w.validatedFragmentSpreads = make(map[string]bool)
		w.walkFragment(child)
	}
}

func (w *Walker) walkOperation(operation *syn.OperationDefinition) {
	w.CurrentOperation = operation
	for _, varDef := range operation.VariableDefinitions {
		varDef.Definition = w.Schema.Types[varDef.Type.Name()]
		for _, v := range w.Observers.variable {
			v(w, varDef)
		}
		if varDef.DefaultValue != nil {
			varDef.DefaultValue.ExpectedType = varDef.Type
			varDef.DefaultValue.Definition = w.Schema.Types[varDef.Type.Name()]
		}
	}

	var def *syn.Definition
	var loc syn.DirectiveLocation
	switch operation.Operation {
	case syn.Query, "":
		def = w.Schema.Query
		loc = syn.LocationQuery
	case syn.Mutation:
		def = w.Schema.Mutation
		loc = syn.LocationMutation
	case syn.Subscription:
		def = w.Schema.Subscription
		loc = syn.LocationSubscription
	}

	for _, varDef := range operation.VariableDefinitions {
		if varDef.DefaultValue != nil {
			w.walkValue(varDef.DefaultValue)
		}
		w.walkDirectives(varDef.Definition, varDef.Directives, syn.LocationVariableDefinition)
	}

	w.walkDirectives(def, operation.Directives, loc)
	w.walkSelectionSet(def, operation.SelectionSet)

	for _, v := range w.Observers.operationVisitor {
		v(w, operation)
	}
	w.CurrentOperation = nil
}

func (w *Walker) walkFragment(it *syn.FragmentDefinition) {
	def := w.Schema.Types[it.TypeCondition.Value]

	it.Definition = def

	w.walkDirectives(def, it.Directives, syn.LocationFragmentDefinition)
	w.walkSelectionSet(def, it.SelectionSet)

	for _, v := range w.Observers.fragment {
		v(w, it)
	}
}

func (w *Walker) walkDirectives(parentDef *syn.Definition, directives []*syn.Directive, location syn.DirectiveLocation) {
	for _, dir := range directives {
		def := w.Schema.Directives[dir.Name]
		dir.Definition = def
		dir.ParentDefinition = parentDef
		dir.Location = location

		for _, arg := range dir.Arguments {
			var argDef *syn.ArgumentDefinition
			if def != nil {
				argDef = def.Arguments.ForName(arg.Name)
			}

			w.walkArgument(argDef, arg)
		}

		for _, v := range w.Observers.directive {
			v(w, dir)
		}
	}

	for _, v := range w.Observers.directiveList {
		v(w, directives)
	}
}

func (w *Walker) walkValue(value *syn.Value) {
	if value.Kind == syn.Variable && w.CurrentOperation != nil {
		value.VariableDefinition = w.CurrentOperation.VariableDefinitions.ForName(value.Raw)
		if value.VariableDefinition != nil {
			value.VariableDefinition.Used = true
		}
	}

	if value.Kind == syn.ObjectValue {
		for _, child := range value.Children {
			if value.Definition != nil {
				fieldDef := value.Definition.Fields.ForName(child.Name)
				if fieldDef != nil {
					child.Value.ExpectedType = fieldDef.Type
					child.Value.Definition = w.Schema.Types[fieldDef.Type.Name()]
				}
			}
			w.walkValue(child.Value)
		}
	}

	if value.Kind == syn.ListValue {
		for _, child := range value.Children {
			if value.ExpectedType != nil && value.ExpectedType.Elem != nil {
				child.Value.ExpectedType = value.ExpectedType.Elem
				child.Value.Definition = value.Definition
			}

			w.walkValue(child.Value)
		}
	}

	for _, v := range w.Observers.value {
		v(w, value)
	}
}

func (w *Walker) walkArgument(argDef *syn.ArgumentDefinition, arg *syn.Argument) {
	if argDef != nil {
		arg.Value.ExpectedType = argDef.Type
		arg.Value.Definition = w.Schema.Types[argDef.Type.Name()]
	}

	w.walkValue(arg.Value)
}

func (w *Walker) walkSelectionSet(parentDef *syn.Definition, it syn.SelectionSet) {
	for _, child := range it {
		w.walkSelection(parentDef, child)
	}
}

func (w *Walker) walkSelection(parentDef *syn.Definition, it syn.Selection) {
	switch it := it.(type) {
	case *syn.Field:
		var def *syn.FieldDefinition
		if it.Name.Value == "__typename" {
			def = &syn.FieldDefinition{
				Name: "__typename",
				Type: syn.NamedType("String", ast.Position{}),
			}
		} else if parentDef != nil {
			def = parentDef.Fields.ForName(it.Name.Value)
		}

		it.Definition = def
		it.ObjectDefinition = parentDef

		var nextParentDef *syn.Definition
		if def != nil {
			nextParentDef = w.Schema.Types[def.Type.Name()]
		}

		for _, arg := range it.Arguments {
			var argDef *syn.ArgumentDefinition
			if def != nil {
				argDef = def.Arguments.ForName(arg.Name)
			}

			w.walkArgument(argDef, arg)
		}

		w.walkDirectives(nextParentDef, it.Directives, syn.LocationField)
		w.walkSelectionSet(nextParentDef, it.SelectionSet)

		for _, v := range w.Observers.field {
			v(w, it)
		}

	case *syn.InlineFragment:
		it.ObjectDefinition = parentDef

		nextParentDef := parentDef
		if it.TypeCondition != "" {
			nextParentDef = w.Schema.Types[it.TypeCondition]
		}

		w.walkDirectives(nextParentDef, it.Directives, syn.LocationInlineFragment)
		w.walkSelectionSet(nextParentDef, it.SelectionSet)

		for _, v := range w.Observers.inlineFragment {
			v(w, it)
		}

	case *syn.FragmentSpread:
		def := w.Document.Fragments.ForName(it.Name)
		it.Definition = def
		it.ObjectDefinition = parentDef

		var nextParentDef *syn.Definition
		if def != nil {
			nextParentDef = w.Schema.Types[def.TypeCondition.Value]
		}

		w.walkDirectives(nextParentDef, it.Directives, syn.LocationFragmentSpread)

		if def != nil && !w.validatedFragmentSpreads[def.Name.Value] {
			// prevent infinite recursion
			w.validatedFragmentSpreads[def.Name.Value] = true
			w.walkSelectionSet(nextParentDef, def.SelectionSet)
		}

		for _, v := range w.Observers.fragmentSpread {
			v(w, it)
		}

	default:
		panic(fmt.Errorf("unsupported %T", it))
	}
}
