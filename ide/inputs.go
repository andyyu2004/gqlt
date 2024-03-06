package ide

import (
	"github.com/movio/gqlt/internal/config"
	"github.com/movio/gqlt/memosa"
)

type sourcesInputQuery struct{}

var _ memosa.Input[Input] = sourcesInputQuery{}

func (sourcesInputQuery) Execute(*memosa.Context, memosa.InputKey) Input {
	panic("memosa will not call input queries")
}

type Input struct {
	Sources map[string]string
}

type schemasInputQuery struct{}

var _ memosa.Input[*config.Schemas] = schemasInputQuery{}

func (schemasInputQuery) Execute(ctx *memosa.Context, key memosa.InputKey) *config.Schemas {
	panic("memosa will not call input queries")
}
