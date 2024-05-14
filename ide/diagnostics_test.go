package ide_test

import (
	"testing"

	"github.com/andyyu2004/expect-test"
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/ide/mapper"
	"github.com/movio/gqlt/internal/annotate"
	"github.com/movio/gqlt/internal/slice"
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

func testDiagnostics(t *testing.T, cases []diagnosticTestCase, minSeverity protocol.DiagnosticSeverity) {
	check := func(name, content string, expectation expect.Expectation) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ide.TestWith(t, content, func(uri string, s ide.Snapshot) {
				mapper := s.Mapper(uri)
				diagnostics := slice.Filter(s.Diagnostics()[uri], func(d protocol.Diagnostic) bool {
					if d.Severity == nil {
						return true
					}

					// Low number is high severity for some reason
					return *d.Severity <= minSeverity
				})

				annotations := slice.Map(diagnostics, toAnnotation(mapper))
				expectation.AssertEqual(t, annotate.Annotate(content, annotations))
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
#        ^ unused variable 'y'
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

		{
			"undefined fragment", `
query { foo { ...Foo } }
`,
			expect.Expect(`
query { foo { ...Foo } }
#                ^^^ fragment 'Foo' not defined
`),
		},

		{
			"unused variable in matches", `
let x = 5
assert 5 matches x
`,
			expect.Expect(`
let x = 5
#   ^ unused variable 'x'
assert 5 matches x
#                ^ unused variable 'x'
`),
		},

		{
			"underscore variables are ignored by unused var", `
let _x = 5
`,
			expect.Expect(`
let _x = 5
`),
		},
	}, protocol.DiagnosticSeverityHint)
}
