package ast

import (
	"fmt"
	"strings"
)

type Errors []Error

func (errs Errors) Error() string {
	s := new(strings.Builder)
	for _, err := range errs {
		fmt.Fprintf(s, "%s\n", err.Error())
	}
	return s.String()
}

type Error struct {
	Position
	Msg string
}

func (err Error) Error() string {
	return fmt.Sprintf("%v: %s", err.Position, err.Message())
}

func (err Error) Pos() Position {
	return err.Position
}

func (err Error) Message() string {
	return err.Msg
}
