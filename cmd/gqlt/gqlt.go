package main

import (
	"log/slog"

	"github.com/andyyu2004/gqlt/ide"
	"github.com/andyyu2004/gqlt/lsp"
)

func main() {
	ide := ide.New()
	server := lsp.New(ide)
	if err := server.RunStdio(); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}
