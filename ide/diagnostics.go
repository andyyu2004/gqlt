package ide

import (
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/andyyu2004/gqlt/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s Snapshot) Diagnostics() map[string][]protocol.Diagnostic {
	sources := s.ide.Sources()
	diagnostics := map[string][]protocol.Diagnostic{}
	for path := range sources {
		diagnostics[path] = s.diagnoseFile(path)
	}
	return diagnostics
}

func (s Snapshot) diagnoseFile(path string) []protocol.Diagnostic {
	root := s.Parse(path)
	mapper := s.Mapper(path)
	diagnostics := []protocol.Diagnostic{}
	if root.Err != nil {
		errs := root.Err.(parser.Errors)
		for _, err := range errs {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range:    posToProto(mapper, err.Position),
				Severity: lib.Ref(protocol.DiagnosticSeverityError),
				Message:  err.Message,
			})
		}
	}
	return diagnostics
}
