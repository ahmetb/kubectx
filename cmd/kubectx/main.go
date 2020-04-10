package main

import (
	"fmt"
	"os"
)

func main() {
	// parse command-line flags
	argv := os.Args[1:]
	fmt.Printf("%#v\n", argv)

	//var op Op
	//op _= parseArgs(argv) // -> DeleteOp RenameOp HelpOp UnrecognizedFlags
}

