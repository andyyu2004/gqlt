package memosa_test

import (
	"testing"

	"github.com/andyyu2004/memosa"
)

type filesInput struct{}

type Files map[string]string

var _ memosa.Input[Files] = filesInput{}

func (filesInput) Execute(*memosa.Context, memosa.InputKey) Files {
	panic("unexpected call to input execute")
}

func TestCalculator(t *testing.T) {
	ctx := memosa.New()
	memosa.Set[filesInput](ctx, Files{})
}
