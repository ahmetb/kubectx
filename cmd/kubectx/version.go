package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// VersionOp show kubectx version.
type VersionOp struct{}

func (_ VersionOp) Run(stdout, _ io.Writer) error {
	return printVersion(stdout)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "Version: %s\nCommit: %s\nDate: %s", version, commit, date)
	return errors.Wrap(err, "write error")
}
