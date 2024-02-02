package ide

import protocol "github.com/tliron/glsp/protocol_3_16"

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
	_ = root
	return nil
}
