package main

import (
	"fmt"
	"io"
)

var (
	version = "v0.0.0+unknown" // populated by goreleaser
)

// VersionOps describes printing version string.
type VersionOp struct{}

func (_ VersionOp) Run(stdout, _ io.Writer) error {
	_, err := fmt.Fprintf(stdout, "%s\n", version)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}
