package main

import (
	"io"
)

type SwitchOp struct{ Target string }

func (s SwitchOp) Run(stdout, stderr io.Writer) error {
	panic("implement me")
}
