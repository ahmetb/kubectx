package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	// parse command-line flags
	var op Op
	op = parseArgs(os.Args[1:])

	if err := op.Run(os.Stdout, os.Stderr); err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, color.RedString("error: ")+format+"\n", args...)
}

func printWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, color.YellowString("warning: ")+format+"\n", args...)
}
