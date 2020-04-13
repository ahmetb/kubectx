package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/printer"
)

type Op interface {
	Run(stdout, stderr io.Writer) error
}

func main() {
	op := parseArgs(os.Args[1:])
	if err := op.Run(os.Stdout, os.Stderr); err != nil {
		printer.Error(os.Stderr, err.Error())

		if _, ok := os.LookupEnv(env.EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(os.Stderr, "[DEBUG] error: %+v\n", err)
		}
		defer os.Exit(1)
	}
}

