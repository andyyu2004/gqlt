package ide

import (
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/internal/typecheck"
	"github.com/andyyu2004/gqlt/memosa/lib"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type diagnostics struct {
	Snapshot
	path        string
	diagnostics []protocol.Diagnostic
}

func (s Snapshot) Diagnostics() map[string][]protocol.Diagnostic {
	sources := s.ide.Sources()
	diagnostics := map[string][]protocol.Diagnostic{}
	for path := range sources {
		diagnostics[path] = s.diagnostics(path)
	}
	return diagnostics
}

func (s Snapshot) diagnostics(path string) []protocol.Diagnostic {
	d := &diagnostics{s, path, []protocol.Diagnostic{}}
	d.diagnose()
	return d.diagnostics
}

func (d *diagnostics) diagnose() {
	d.syntax()
	d.typecheck()
}

func (d *diagnostics) syntax() {
	root := d.Parse(d.path)
	mapper := d.Mapper(d.path)
	if root.Err != nil {
		if errs, ok := root.Err.(ast.Errors); ok {
			for _, err := range errs {
				d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
					Range:    posToProto(mapper, err.Position),
					Severity: lib.Ref(protocol.DiagnosticSeverityError),
					Message:  err.Message(),
				})
			}
		}
	}
}

func (d *diagnostics) typecheck() {
	root := d.Parse(d.path)
	tcx := typecheck.New()
	_, err := tcx.Check(root.Ast)
	if err != nil {
		errs := err.(typecheck.Errors)
		for _, err := range errs {
			d.diagnostics = append(d.diagnostics, protocol.Diagnostic{
				Range:    posToProto(d.Mapper(d.path), err.Pos),
				Severity: lib.Ref(protocol.DiagnosticSeverityError),
				Message:  err.Msg,
			})
		}
	}
}
