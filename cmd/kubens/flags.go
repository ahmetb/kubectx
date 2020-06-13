package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

// UnsupportedOp indicates an unsupported flag.
type UnsupportedOp struct{ Err error }

func (op UnsupportedOp) Run(_, _ io.Writer) error {
	return op.Err
}

// parseArgs looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func parseArgs(argv []string) Op {
	if len(argv) == 0 {
		if cmdutil.IsInteractiveMode(os.Stdout) {
			return InteractiveSwitchOp{SelfCmd: os.Args[0]}
		}
		return ListOp{}
	}

	if len(argv) == 1 {
		v := argv[0]
		if v == "--help" || v == "-h" {
			return HelpOp{}
		}
		if v == "--current" || v == "-c" {
			return CurrentOp{}
		}
		if v == "--version" || v == "-v" {
			return VersionOp{}
		}
		if strings.HasPrefix(v, "-") && v != "-" {
			return UnsupportedOp{Err: fmt.Errorf("unsupported option '%s'", v)}
		}
		return SwitchOp{Target: argv[0]}
	}
	return UnsupportedOp{Err: fmt.Errorf("too many arguments")}
}
