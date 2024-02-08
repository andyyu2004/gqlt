package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/ide"
	"github.com/andyyu2004/gqlt/ide/mapper"
	"github.com/andyyu2004/gqlt/internal/annotate"
	"github.com/andyyu2004/gqlt/internal/slice"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func toAnnotation(mapper *mapper.Mapper) func(diagnostic protocol.Diagnostic) annotate.Annotation {
	return func(diag protocol.Diagnostic) annotate.Annotation {
		return ann{diag, mapper}
	}
}

type ann struct {
	protocol.Diagnostic
	mapper *mapper.Mapper
}

func (a ann) Message() string {
	return a.Diagnostic.Message
}

func (a ann) Pos() ast.Position {
	start, err := a.mapper.ByteOffset(int(a.Range.Start.Line), int(a.Range.Start.Character))
	if err != nil {
		panic(err)
	}

	end, err := a.mapper.ByteOffset(int(a.Range.End.Line), int(a.Range.End.Character))
	if err != nil {
		panic(err)
	}

	return ast.Position{
		Start:  start,
		End:    end,
		Line:   int(a.Range.Start.Line + 1),
		Column: int(a.Range.Start.Character + 1),
	}
}

type diagnosticTestCase struct {
	name        string
	content     string
	expectation expect.Expectation
}

func testDiagnostics(t *testing.T, cases ...diagnosticTestCase) {
	check := func(name, content string, expectation expect.Expectation) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ide.TestWith(t, content, func(uri string, s ide.Snapshot) {
				mapper := s.Mapper(uri)
				diagnostics := slice.Map(s.Diagnostics()[uri], toAnnotation(mapper))

				expectation.AssertEqual(t, annotate.Annotate(content, diagnostics))
			})
		})
	}

	for _, c := range cases {
		check(c.name, c.content, c.expectation)
	}
}

func TestDiagnostics(t *testing.T) {
	testDiagnostics(t, []diagnosticTestCase{
		{
			"syntax",
			"let x == 5",
			expect.Expect(`let x == 5
#     ^ expected '=', found '=='`),
		},

		{
			"object pattern alias hides original name", `
let { x: y } = { x: 5 }
print(x)
`,
			expect.Expect(`
let { x: y } = { x: 5 }
print(x)
#     ^ unbound name 'x'
`),
		},

		{
			"object pattern nested pattern names", `
let { x: { y: z } } = { x: { y: 1 } }
print(x)
print(y)
print(z)
`,
			expect.Expect(`
let { x: { y: z } } = { x: { y: 1 } }
print(x)
#     ^ unbound name 'x'
print(y)
#     ^ unbound name 'y'
print(z)
`),
		},

		{
			"object pattern nested pattern names", `
let { errors: [{ message }] } = try query { foo }
print(errors)
print(message)
`,
			expect.Expect(`
let { errors: [{ message }] } = try query { foo }
#                                           ^^^ field 'foo' of type '{ id: ID, any: Any, string: String, int: Int, float: Float, boolean: Boolean }' must have a selection of subfields
print(errors)
#     ^^^^^^ unbound name 'errors'
print(message)
`),
		},
	}...)
}
