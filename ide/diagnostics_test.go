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
			ide.TestWith(t, content, func(path string, s ide.Snapshot) {
				mapper := s.Mapper(path)
				diagnostics := slice.Map(s.Diagnostics()[path], toAnnotation(mapper))

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
	}...)
	//	check("set namespace", `
	//
	// set namespace = "foo"
	// set namespace ["a", "b", "c"]
	// set namespace false`, expect.Expect(`
	// set namespace = "foo"
	// set namespace ["a", "b", "c"]
	// set namespace false
	// #^^^^^^^^^^^^^^^^^^^ expected string or list of strings as value for "namespace", found bool
	// `))
}
