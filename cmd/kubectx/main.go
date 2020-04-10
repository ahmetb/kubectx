package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// parse command-line flags
	var op Op
	op = parseArgs(os.Args[1:])

	switch v := op.(type) {
	case ListOp:
		// TODO implement
		panic("not implemented")
	case SwitchOp:
		// TODO implement
		panic("not implemented")
	case UnknownOp:
		fmt.Printf("error: unsupported operation: %s\n", strings.Join(v.Args, " "))
		// TODO print --help string
		os.Exit(1)
	default:
		fmt.Printf("internal error: operation type %T not handled", op)
	}
}
