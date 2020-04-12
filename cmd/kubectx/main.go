package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

func main() {
	// parse command-line flags
	var op Op
	op = parseArgs(os.Args[1:])

	if err := op.Run(os.Stdout, os.Stderr); err != nil {
		printError(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func printError(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, color.RedString("error: ")+format+"\n", args...)
}

func printWarning(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.YellowString("warning: ")+format+"\n", args...)
}

func printSuccess(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.GreenString(fmt.Sprintf(format+"\n", args...)))
}
