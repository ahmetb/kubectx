package main

import "io"

type CurrentOp struct{}

func (c CurrentOp) Run(stdout, stderr io.Writer) error {
	panic("implement me")
}
