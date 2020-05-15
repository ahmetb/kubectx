package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/printer"
	"github.com/fatih/color"
)

type Op interface {
	Run(stdout, stderr io.Writer) error
}

func main() {
	cmdutil.PrintDeprecatedEnvWarnings(color.Error, os.Environ())

	op := parseArgs(os.Args[1:])
	if err := op.Run(color.Output, color.Error); err != nil {
		printer.Error(color.Error, err.Error())

		if _, ok := os.LookupEnv(env.EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(color.Error, "[DEBUG] error: %+v\n", err)
		}
		defer os.Exit(1)
	}
}
