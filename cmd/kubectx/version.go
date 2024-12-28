package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	version = "v0.0.0+unknown" // populated by goreleaser
)

// VersionOp describes printing version string.
type VersionOp struct{}

func (_ VersionOp) Run(stdout, _ io.Writer) error {
	_, err := fmt.Fprintf(stdout, "%s\n", version)
	return errors.Wrap(err, "write error")
}
