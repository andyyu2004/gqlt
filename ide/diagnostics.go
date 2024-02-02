package ide

import protocol "github.com/tliron/glsp/protocol_3_16"

func (s Snapshot) Diagnostics() map[string][]protocol.Diagnostic {
	return map[string][]protocol.Diagnostic{}
}
