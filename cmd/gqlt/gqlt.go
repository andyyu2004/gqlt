package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/movio/gqlt"
	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/lsp"
)

var args struct {
	Run   runCmd   `cmd:"run" help:"Run a query against a GraphQL server."`
	Serve serveCmd `cmd:"serve" default:"1" help:"Run the gqlt lsp server using stdio"`
}

type runCmd struct {
	URL  string `arg:"" help:"URL of the GraphQL server to query"`
	File string `arg:"" help:"Path to the file to run (use - for stdin)"`
}

func (r *runCmd) Run() error {
	file := r.File
	if r.File == "-" {
		file = "/dev/stdin"
	}

	executor := gqlt.New()
	ctx := context.Background()

	return executor.RunFile(ctx, gqlt.HTTPClient{Client: http.DefaultClient, URL: r.URL}, file)
}

type serveCmd struct{}

func (s *serveCmd) Run() error {
	ide := ide.New()
	server := lsp.New(ide)
	return server.RunStdio()
}

func main() {
	ctx := kong.Parse(&args)
	if err := ctx.Run(); err != nil {
		slog.Error("failed to run", "error", err)
	}
}
