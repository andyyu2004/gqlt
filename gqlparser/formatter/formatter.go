package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/andyyu2004/gqlt/syn"
)

type Formatter interface {
	FormatSchema(schema *syn.Schema)
	FormatSchemaDocument(doc *syn.SchemaDocument)
	FormatQueryDocument(doc *syn.QueryDocument)
}

//nolint:revive // Ignore "stuttering" name format.FormatterOption
type FormatterOption func(*formatter)

// WithIndent uses the given string for indenting block bodies in the output,
// instead of the default, `"\t"`.
func WithIndent(indent string) FormatterOption {
	return func(f *formatter) {
		f.indent = indent
	}
}

// WithComments includes comments from the source/AST in the formatted output.
func WithComments() FormatterOption {
	return func(f *formatter) {
		f.emitComments = true
	}
}

func NewFormatter(w io.Writer, options ...FormatterOption) Formatter {
	f := &formatter{
		indent: "\t",
		writer: w,
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type formatter struct {
	writer io.Writer

	indent       string
	indentSize   int
	emitBuiltin  bool
	emitComments bool

	padNext  bool
	lineHead bool
}

func (f *formatter) writeString(s string) {
	_, _ = f.writer.Write([]byte(s))
}

func (f *formatter) writeIndent() *formatter {
	if f.lineHead {
		f.writeString(strings.Repeat(f.indent, f.indentSize))
	}
	f.lineHead = false
	f.padNext = false

	return f
}

func (f *formatter) WriteNewline() *formatter {
	f.writeString("\n")
	f.lineHead = true
	f.padNext = false

	return f
}

func (f *formatter) WriteWord(word string) *formatter {
	if f.lineHead {
		f.writeIndent()
	}
	if f.padNext {
		f.writeString(" ")
	}
	f.writeString(strings.TrimSpace(word))
	f.padNext = true

	return f
}

func (f *formatter) WriteString(s string) *formatter {
	if f.lineHead {
		f.writeIndent()
	}
	if f.padNext {
		f.writeString(" ")
	}
	f.writeString(s)
	f.padNext = false

	return f
}

func (f *formatter) WriteDescription(s string) *formatter {
	if s == "" {
		return f
	}

	f.WriteString(`"""`)
	if ss := strings.Split(s, "\n"); len(ss) > 1 {
		f.WriteNewline()
		for _, s := range ss {
			f.WriteString(s).WriteNewline()
		}
	} else {
		f.WriteString(s)
	}

	f.WriteString(`"""`).WriteNewline()

	return f
}

func (f *formatter) IncrementIndent() {
	f.indentSize++
}

func (f *formatter) DecrementIndent() {
	f.indentSize--
}

func (f *formatter) NoPadding() *formatter {
	f.padNext = false

	return f
}

func (f *formatter) NeedPadding() *formatter {
	f.padNext = true

	return f
}

func (f *formatter) FormatSchema(schema *syn.Schema) {
	if schema == nil {
		return
	}

	f.FormatCommentGroup(schema.Comment)

	var inSchema bool
	startSchema := func() {
		if !inSchema {
			inSchema = true

			f.WriteWord("schema").WriteString("{").WriteNewline()
			f.IncrementIndent()
		}
	}
	if schema.Query != nil && schema.Query.Name != "Query" {
		startSchema()
		f.WriteWord("query").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Query.Name).WriteNewline()
	}
	if schema.Mutation != nil && schema.Mutation.Name != "Mutation" {
		startSchema()
		f.WriteWord("mutation").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Mutation.Name).WriteNewline()
	}
	if schema.Subscription != nil && schema.Subscription.Name != "Subscription" {
		startSchema()
		f.WriteWord("subscription").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Subscription.Name).WriteNewline()
	}
	if inSchema {
		f.DecrementIndent()
		f.WriteString("}").WriteNewline()
	}

	directiveNames := make([]string, 0, len(schema.Directives))
	for name := range schema.Directives {
		directiveNames = append(directiveNames, name)
	}
	sort.Strings(directiveNames)
	for _, name := range directiveNames {
		f.FormatDirectiveDefinition(schema.Directives[name])
	}

	typeNames := make([]string, 0, len(schema.Types))
	for name := range schema.Types {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)
	for _, name := range typeNames {
		f.FormatDefinition(schema.Types[name], false)
	}
}

func (f *formatter) FormatSchemaDocument(doc *syn.SchemaDocument) {
	// TODO emit by position based order

	if doc == nil {
		return
	}

	f.FormatSchemaDefinitionList(doc.Schema, false)
	f.FormatSchemaDefinitionList(doc.SchemaExtension, true)

	f.FormatDirectiveDefinitionList(doc.Directives)

	f.FormatDefinitionList(doc.Definitions, false)
	f.FormatDefinitionList(doc.Extensions, true)

	// doc.Comment is end of file comment, so emit last
	f.FormatCommentGroup(doc.Comment)
}

func (f *formatter) FormatQueryDocument(doc *syn.QueryDocument) {
	// TODO emit by position based order

	if doc == nil {
		return
	}

	f.FormatCommentGroup(doc.Comment)

	f.FormatOperationList(doc.Operations)
	f.FormatFragmentDefinitionList(doc.Fragments)
}

func (f *formatter) FormatSchemaDefinitionList(lists syn.SchemaDefinitionList, extension bool) {
	if len(lists) == 0 {
		return
	}

	var (
		beforeDescComment      = new(syn.CommentGroup)
		afterDescComment       = new(syn.CommentGroup)
		endOfDefinitionComment = new(syn.CommentGroup)
		description            string
	)

	for _, def := range lists {
		if def.BeforeDescriptionComment != nil {
			beforeDescComment.List = append(beforeDescComment.List, def.BeforeDescriptionComment.List...)
		}
		if def.AfterDescriptionComment != nil {
			afterDescComment.List = append(afterDescComment.List, def.AfterDescriptionComment.List...)
		}
		if def.EndOfDefinitionComment != nil {
			endOfDefinitionComment.List = append(endOfDefinitionComment.List, def.EndOfDefinitionComment.List...)
		}
		description += def.Description
	}

	f.FormatCommentGroup(beforeDescComment)
	f.WriteDescription(description)
	f.FormatCommentGroup(afterDescComment)

	if extension {
		f.WriteWord("extend")
	}
	f.WriteWord("schema").WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, def := range lists {
		f.FormatSchemaDefinition(def)
	}

	f.FormatCommentGroup(endOfDefinitionComment)

	f.DecrementIndent()
	f.WriteString("}").WriteNewline()
}

func (f *formatter) FormatSchemaDefinition(def *syn.SchemaDefinition) {
	f.FormatDirectiveList(def.Directives)

	f.FormatOperationTypeDefinitionList(def.OperationTypes)
}

func (f *formatter) FormatOperationTypeDefinitionList(lists syn.OperationTypeDefinitionList) {
	for _, def := range lists {
		f.FormatOperationTypeDefinition(def)
	}
}

func (f *formatter) FormatOperationTypeDefinition(def *syn.OperationTypeDefinition) {
	f.FormatCommentGroup(def.Comment)
	f.WriteWord(string(def.Operation)).NoPadding().WriteString(":").NeedPadding()
	f.WriteWord(def.Type)
	f.WriteNewline()
}

func (f *formatter) FormatFieldList(fieldList syn.FieldList, endOfDefComment *syn.CommentGroup) {
	if len(fieldList) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, field := range fieldList {
		f.FormatFieldDefinition(field)
	}

	f.FormatCommentGroup(endOfDefComment)

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatFieldDefinition(field *syn.FieldDefinition) {
	if !f.emitBuiltin && strings.HasPrefix(field.Name, "__") {
		return
	}

	f.FormatCommentGroup(field.BeforeDescriptionComment)

	f.WriteDescription(field.Description)

	f.FormatCommentGroup(field.AfterDescriptionComment)

	f.WriteWord(field.Name).NoPadding()
	f.FormatArgumentDefinitionList(field.Arguments)
	f.NoPadding().WriteString(":").NeedPadding()
	f.FormatType(field.Type)

	if field.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(field.DefaultValue)
	}

	f.FormatDirectiveList(field.Directives)

	f.WriteNewline()
}

func (f *formatter) FormatArgumentDefinitionList(lists syn.ArgumentDefinitionList) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("(")
	for idx, arg := range lists {
		f.FormatArgumentDefinition(arg)

		// Skip emitting (insignificant) comma in case it is the
		// last argument, or we printed a new line in its definition.
		if idx != len(lists)-1 && arg.Description == "" {
			f.NoPadding().WriteWord(",")
		}
	}
	f.NoPadding().WriteString(")").NeedPadding()
}

func (f *formatter) FormatArgumentDefinition(def *syn.ArgumentDefinition) {
	f.FormatCommentGroup(def.BeforeDescriptionComment)

	if def.Description != "" {
		f.WriteNewline().IncrementIndent()
		f.WriteDescription(def.Description)
	}

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord(def.Name).NoPadding().WriteString(":").NeedPadding()
	f.FormatType(def.Type)

	if def.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(def.DefaultValue)
	}

	f.NeedPadding().FormatDirectiveList(def.Directives)

	if def.Description != "" {
		f.DecrementIndent()
		f.WriteNewline()
	}
}

func (f *formatter) FormatDirectiveLocation(location syn.DirectiveLocation) {
	f.WriteWord(string(location))
}

func (f *formatter) FormatDirectiveDefinitionList(lists syn.DirectiveDefinitionList) {
	if len(lists) == 0 {
		return
	}

	for _, dec := range lists {
		f.FormatDirectiveDefinition(dec)
	}
}

func (f *formatter) FormatDirectiveDefinition(def *syn.DirectiveDefinition) {
	if !f.emitBuiltin {
		if def.Position.Src.BuiltIn {
			return
		}
	}

	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord("directive").WriteString("@").WriteWord(def.Name)

	if len(def.Arguments) != 0 {
		f.NoPadding()
		f.FormatArgumentDefinitionList(def.Arguments)
	}

	if def.IsRepeatable {
		f.WriteWord("repeatable")
	}

	if len(def.Locations) != 0 {
		f.WriteWord("on")

		for idx, dirLoc := range def.Locations {
			f.FormatDirectiveLocation(dirLoc)

			if idx != len(def.Locations)-1 {
				f.WriteWord("|")
			}
		}
	}

	f.WriteNewline()
}

func (f *formatter) FormatDefinitionList(lists syn.DefinitionList, extend bool) {
	if len(lists) == 0 {
		return
	}

	for _, dec := range lists {
		f.FormatDefinition(dec, extend)
	}
}

func (f *formatter) FormatDefinition(def *syn.Definition, extend bool) {
	if !f.emitBuiltin && def.BuiltIn {
		return
	}

	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	if extend {
		f.WriteWord("extend")
	}

	switch def.Kind {
	case syn.Scalar:
		f.WriteWord("scalar").WriteWord(def.Name)

	case syn.Object:
		f.WriteWord("type").WriteWord(def.Name)

	case syn.Interface:
		f.WriteWord("interface").WriteWord(def.Name)

	case syn.Union:
		f.WriteWord("union").WriteWord(def.Name)

	case syn.Enum:
		f.WriteWord("enum").WriteWord(def.Name)

	case syn.InputObject:
		f.WriteWord("input").WriteWord(def.Name)
	}

	if len(def.Interfaces) != 0 {
		f.WriteWord("implements").WriteWord(strings.Join(def.Interfaces, " & "))
	}

	f.FormatDirectiveList(def.Directives)

	if len(def.Types) != 0 {
		f.WriteWord("=").WriteWord(strings.Join(def.Types, " | "))
	}

	f.FormatFieldList(def.Fields, def.EndOfDefinitionComment)

	f.FormatEnumValueList(def.EnumValues, def.EndOfDefinitionComment)

	f.WriteNewline()
}

func (f *formatter) FormatEnumValueList(lists syn.EnumValueList, endOfDefComment *syn.CommentGroup) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, v := range lists {
		f.FormatEnumValueDefinition(v)
	}

	f.FormatCommentGroup(endOfDefComment)

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatEnumValueDefinition(def *syn.EnumValueDefinition) {
	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord(def.Name)
	f.FormatDirectiveList(def.Directives)

	f.WriteNewline()
}

func (f *formatter) FormatOperationList(lists syn.OperationList) {
	for _, def := range lists {
		f.FormatOperationDefinition(def)
	}
}

func (f *formatter) FormatOperationDefinition(def *syn.OperationDefinition) {
	f.FormatCommentGroup(def.Comment)

	f.WriteWord(string(def.Operation))
	if def.Name.Value != "" {
		f.WriteWord(def.Name.Value)
	}
	f.FormatVariableDefinitionList(def.VariableDefinitions)
	f.FormatDirectiveList(def.Directives)

	if len(def.SelectionSet) != 0 {
		f.FormatSelectionSet(def.SelectionSet)
		f.WriteNewline()
	}
}

func (f *formatter) FormatDirectiveList(lists syn.DirectiveList) {
	if len(lists) == 0 {
		return
	}

	for _, dir := range lists {
		f.FormatDirective(dir)
	}
}

func (f *formatter) FormatDirective(dir *syn.Directive) {
	f.WriteString("@").WriteWord(dir.Name)
	f.FormatArgumentList(dir.Arguments)
}

func (f *formatter) FormatArgumentList(lists syn.ArgumentList) {
	if len(lists) == 0 {
		return
	}
	f.NoPadding().WriteString("(")
	for idx, arg := range lists {
		f.FormatArgument(arg)

		if idx != len(lists)-1 {
			f.NoPadding().WriteWord(",")
		}
	}
	f.WriteString(")").NeedPadding()
}

func (f *formatter) FormatArgument(arg *syn.Argument) {
	f.FormatCommentGroup(arg.Comment)

	f.WriteWord(arg.Name).NoPadding().WriteString(":").NeedPadding()
	f.WriteString(arg.Value.String())
}

func (f *formatter) FormatFragmentDefinitionList(lists syn.FragmentDefinitionList) {
	for _, def := range lists {
		f.FormatFragmentDefinition(def)
	}
}

func (f *formatter) FormatFragmentDefinition(def *syn.FragmentDefinition) {
	f.FormatCommentGroup(def.Comment)

	f.WriteWord("fragment").WriteWord(def.Name.Value)
	f.FormatVariableDefinitionList(def.VariableDefinition)
	f.WriteWord("on").WriteWord(def.TypeCondition.Value)
	f.FormatDirectiveList(def.Directives)

	if len(def.SelectionSet) != 0 {
		f.FormatSelectionSet(def.SelectionSet)
		f.WriteNewline()
	}
}

func (f *formatter) FormatVariableDefinitionList(lists syn.VariableDefinitionList) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("(")
	for idx, def := range lists {
		f.FormatVariableDefinition(def)

		if idx != len(lists)-1 {
			f.NoPadding().WriteWord(",")
		}
	}
	f.NoPadding().WriteString(")").NeedPadding()
}

func (f *formatter) FormatVariableDefinition(def *syn.VariableDefinition) {
	f.FormatCommentGroup(def.Comment)

	f.WriteString("$").WriteWord(def.Variable).NoPadding().WriteString(":").NeedPadding()
	f.FormatType(def.Type)

	if def.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(def.DefaultValue)
	}

	// TODO https://github.com/andyyu2004/gqlt/gqlparser/issues/102
	//   VariableDefinition : Variable : Type DefaultValue? Directives[Const]?
}

func (f *formatter) FormatSelectionSet(sets syn.SelectionSet) {
	if len(sets) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, sel := range sets {
		f.FormatSelection(sel)
	}

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatSelection(selection syn.Selection) {
	switch v := selection.(type) {
	case *syn.Field:
		f.FormatField(v)

	case *syn.FragmentSpread:
		f.FormatFragmentSpread(v)

	case *syn.InlineFragment:
		f.FormatInlineFragment(v)

	default:
		panic(fmt.Errorf("unknown Selection type: %T", selection))
	}

	f.WriteNewline()
}

func (f *formatter) FormatField(field *syn.Field) {
	f.FormatCommentGroup(field.Comment)

	if field.Alias.Value != "" && field.Alias.Value != field.Name.Value {
		f.WriteWord(field.Alias.Value).NoPadding().WriteString(":").NeedPadding()
	}
	f.WriteWord(field.Name.Value)

	if len(field.Arguments) != 0 {
		f.NoPadding()
		f.FormatArgumentList(field.Arguments)
		f.NeedPadding()
	}

	f.FormatDirectiveList(field.Directives)

	f.FormatSelectionSet(field.SelectionSet)
}

func (f *formatter) FormatFragmentSpread(spread *syn.FragmentSpread) {
	f.FormatCommentGroup(spread.Comment)

	f.WriteWord("...").WriteWord(spread.Name.Value)

	f.FormatDirectiveList(spread.Directives)
}

func (f *formatter) FormatInlineFragment(inline *syn.InlineFragment) {
	f.FormatCommentGroup(inline.Comment)

	f.WriteWord("...")
	if inline.TypeCondition.Value != "" {
		f.WriteWord("on").WriteWord(inline.TypeCondition.Value)
	}

	f.FormatDirectiveList(inline.Directives)

	f.FormatSelectionSet(inline.SelectionSet)
}

func (f *formatter) FormatType(t *syn.Type) {
	f.WriteWord(t.String())
}

func (f *formatter) FormatValue(value *syn.Value) {
	f.FormatCommentGroup(value.Comment)

	f.WriteString(value.String())
}

func (f *formatter) FormatCommentGroup(group *syn.CommentGroup) {
	if !f.emitComments || group == nil {
		return
	}
	for _, comment := range group.List {
		f.FormatComment(comment)
	}
}

func (f *formatter) FormatComment(comment *syn.Comment) {
	if !f.emitComments || comment == nil {
		return
	}
	f.WriteString("#").WriteString(comment.Text()).WriteNewline()
}
