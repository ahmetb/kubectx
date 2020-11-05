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
	var isForceSwitch bool
	isForceSwitch = false
	ArgsLen := len(argv)
	if ArgsLen == 1 || ArgsLen == 2 {
		v := argv[0]
		if v == "--help" || v == "-h" {
			return HelpOp{}
		}
		if v == "--current" || v == "-c" {
			return CurrentOp{}
		}
		if v == "--force" || v == "-f" {
			isForceSwitch = true
		}
		if strings.HasPrefix(v, "-") && v != "-" && !isForceSwitch {
			return UnsupportedOp{Err: fmt.Errorf("unsupported option '%s'", v)}
		}
		if isForceSwitch {
			if ArgsLen == 1 {
				return UnsupportedOp{Err: fmt.Errorf("force option needs namespace name")}
			}
			return SwitchOp{Target: argv[1], IsForce: isForceSwitch}
		}
		return SwitchOp{Target: argv[0], IsForce: isForceSwitch}
	}
	return UnsupportedOp{Err: fmt.Errorf("too many arguments")}
}
