package ide

import (
	"github.com/andyyu2004/gqlt/memosa"
	"github.com/andyyu2004/gqlt/syn"
)

type sourcesInputQuery struct{}

var _ memosa.Input[Input] = sourcesInputQuery{}

func (sourcesInputQuery) Execute(*memosa.Context, memosa.InputKey) Input {
	panic("memosa will not call input queries")
}

type Input struct {
	Sources map[string]string
}

type schemaInputQuery struct{}

var _ memosa.Input[*syn.Schema] = schemaInputQuery{}

func (schemaInputQuery) Execute(ctx *memosa.Context, key memosa.InputKey) *syn.Schema {
	panic("memosa will not call input queries")
}
