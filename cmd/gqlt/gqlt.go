package main

import (
	"log/slog"

	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/lsp"
)

func main() {
	ide := ide.New()
	server := lsp.New(ide)
	if err := server.RunStdio(); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}
