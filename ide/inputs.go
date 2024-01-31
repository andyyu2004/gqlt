package ide

import (
	"github.com/andyyu2004/gqlt/memosa"
)

type inputQuery struct{}

var _ memosa.Input[Input] = inputQuery{}

func (inputQuery) Execute(*memosa.Context, memosa.InputKey) Input {
	panic("memosa will not call input queries")
}

type Input struct {
	Sources map[string]string
}
