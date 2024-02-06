package lsp

import (
	"errors"
	"fmt"
	"net/url"
	"slices"

	"github.com/andyyu2004/gqlt/ide"
	"github.com/andyyu2004/gqlt/internal/config"
	"github.com/andyyu2004/gqlt/internal/slice"
	"github.com/andyyu2004/gqlt/memosa/lib"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

var semanticTokensLegend = protocol.SemanticTokensLegend{
	TokenTypes: []string{
		string(protocol.SemanticTokenTypeKeyword),
		string(protocol.SemanticTokenTypeProperty),
		string(protocol.SemanticTokenTypeVariable),
		string(protocol.SemanticTokenTypeString),
		string(protocol.SemanticTokenTypeNumber),
		string(protocol.SemanticTokenTypeOperator),
		string(protocol.SemanticTokenTypeType),
		string(protocol.SemanticTokenTypeParameter),
	},
}

func New(ide *ide.IDE) *server.Server {
	ls := &ls{ide}
	handler := &protocol.Handler{
		Initialize:                     ls.initialize,
		Initialized:                    ls.initialized,
		TextDocumentDefinition:         ls.definition,
		TextDocumentDidChange:          ls.onChange,
		TextDocumentDidOpen:            ls.onOpen,
		TextDocumentSemanticTokensFull: ls.semanticTokens,
		TextDocumentHover:              ls.hover,
		SetTrace:                       ls.setTrace,
	}
	server := server.NewServer(handler, "gqlt", false)
	server.Log.Infof("gqlt language server")
	return server
}

type ls struct {
	ide *ide.IDE
}

func (s *ls) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	protocol.SetTraceValue(protocol.TraceValueVerbose)

	workspaces := slice.Map(params.WorkspaceFolders, func(f protocol.WorkspaceFolder) string {
		return lib.Must(url.Parse(f.URI)).Path
	})

	// TODO need to add a lsp file watcher to manage changes to the schema files and config file
	schema, err := config.LoadSchema(workspaces...)
	if err != nil {
		// don't return, just continue without loading a schema
		_ = protocol.Trace(ctx, protocol.MessageTypeError, fmt.Sprintf("failed to load config: %v", err))
	}

	// the input must always be set, nil or not
	s.ide.SetSchema(schema)

	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: lib.Ref(true),
				Change:    lib.Ref(protocol.TextDocumentSyncKindFull),
			},
			SemanticTokensProvider: protocol.SemanticTokensOptions{
				Legend: semanticTokensLegend,
				Full:   true,
			},
			HoverProvider:      protocol.HoverOptions{},
			DefinitionProvider: true,
			Workspace: &protocol.ServerCapabilitiesWorkspace{
				WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{
					Supported: lib.Ref(true),
				},
			},
		},
		ServerInfo: &protocol.InitializeResultServerInfo{Name: "gqlt"},
	}, nil
}

func (s *ls) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (s *ls) onOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.ide.Apply(ide.Changes{
		ide.SetFileContent{
			Path:    params.TextDocument.URI,
			Content: params.TextDocument.Text,
		},
	})
	return nil
}

func (s *ls) onChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	changes := ide.Changes{}
	for _, change := range params.ContentChanges {
		switch change := change.(type) {
		case protocol.TextDocumentContentChangeEvent:
			return errors.New("expected full update events only")
		case protocol.TextDocumentContentChangeEventWhole:
			changes = append(changes, ide.SetFileContent{
				Path:    params.TextDocument.URI,
				Content: change.Text,
			})
		}
	}

	s.ide.Apply(changes)

	s.publishDiagnostics(ctx)

	return nil
}

type logger struct {
	ctx *glsp.Context
}

func (l logger) Debugf(format string, args ...any) {
	_ = protocol.Trace(l.ctx, protocol.MessageTypeLog, fmt.Sprintf(format, args...))
}

func (l logger) Infof(format string, args ...any) {
	_ = protocol.Trace(l.ctx, protocol.MessageTypeInfo, fmt.Sprintf(format, args...))
}

func (l logger) Warnf(format string, args ...any) {
	_ = protocol.Trace(l.ctx, protocol.MessageTypeWarning, fmt.Sprintf(format, args...))
}

func (l logger) Errorf(format string, args ...any) {
	_ = protocol.Trace(l.ctx, protocol.MessageTypeError, fmt.Sprintf(format, args...))
}

func (l *ls) publishDiagnostics(ctx *glsp.Context) {
	log := logger{ctx}
	diagnostics, err := ide.WithSnapshot(l.ide, log, func(s ide.Snapshot) map[string][]protocol.Diagnostic {
		return s.Diagnostics()
	})
	if err != nil {
		log.Errorf("failed to compute diagnostics: %v", err)
		return
	}

	for uri, diags := range diagnostics {
		ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diags,
		})
	}
}

func (l *ls) hover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	return ide.WithSnapshot(l.ide, logger{ctx}, func(s ide.Snapshot) *protocol.Hover {
		return s.Hover(params.TextDocument.URI, params.Position)
	})
}

func (l *ls) semanticTokens(ctx *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	type SemanticToken struct {
		deltaLine            uint32
		deltaStart           uint32
		length               uint32
		tokenType            uint32
		tokenModifiersBitset uint32
	}

	highlights, err := ide.WithSnapshot(l.ide, logger{ctx}, func(s ide.Snapshot) []ide.Highlight {
		return s.Highlight(params.TextDocument.URI)
	})
	if err != nil {
		return nil, err
	}

	tokens := []SemanticToken{}
	for i, hl := range highlights {
		var deltaLine uint32
		// adjust for 1-indexing to 0-indexing
		deltaStart := uint32(hl.Pos.Column - 1)
		if i == 0 {
			deltaLine = uint32(hl.Pos.Line - 1)
		} else {
			deltaLine = uint32(hl.Pos.Line - highlights[i-1].Pos.Line)
			if hl.Pos.Line == highlights[i-1].Pos.Line {
				deltaStart = uint32(hl.Pos.Column - highlights[i-1].Pos.Column)
			}
		}

		tokenType := slices.IndexFunc(semanticTokensLegend.TokenTypes, func(t string) bool { return t == string(hl.TokenKind) })
		lib.Assert(tokenType != -1)
		tokens = append(tokens, SemanticToken{
			deltaLine:            deltaLine,
			deltaStart:           deltaStart,
			length:               uint32(hl.Pos.End - hl.Pos.Start),
			tokenType:            uint32(tokenType),
			tokenModifiersBitset: 0,
		})
	}

	data := make([]uint32, 0, len(tokens)*5)
	for _, token := range tokens {
		data = append(data, token.deltaLine, token.deltaStart, token.length, token.tokenType, token.tokenModifiersBitset)
	}

	return &protocol.SemanticTokens{Data: data}, nil
}

func (l *ls) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (l *ls) definition(ctx *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	return ide.WithSnapshot(l.ide, logger{ctx}, func(s ide.Snapshot) []protocol.Location {
		return s.Definition(params.TextDocument.URI, params.Position)
	})
}
