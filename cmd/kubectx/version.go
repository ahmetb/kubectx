package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// VersionOp show kubectx version.
type VersionOp struct{}

func (_ VersionOp) Run(stdout, _ io.Writer) error {
	return printVersion(stdout)
}

func printVersion(out io.Writer) error {
	version := "v0.9.0"
	_, err := fmt.Fprintf(out, "%s\n", version)
	return errors.Wrap(err, "write error")
}
