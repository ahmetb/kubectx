package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

func main() {
	op := parseArgs(os.Args[1:])
	if err := op.Run(os.Stdout, os.Stderr); err != nil {
		printError(os.Stderr, err.Error())

		if _, ok := os.LookupEnv(EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(os.Stderr, "[DEBUG] error: %+v\n", err)
		}
		defer os.Exit(1)
	}
}

func printError(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.RedString("error: ")+format+"\n", args...)
}

func printWarning(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.YellowString("warning: ")+format+"\n", args...)
}

func printSuccess(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.GreenString(fmt.Sprintf(format+"\n", args...)))
}
