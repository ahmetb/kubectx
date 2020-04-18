package main

import (
	"io"
)

type ListOp struct{}

func (op ListOp) Run(stdout, stderr io.Writer) error {
	panic("implement me")
}
