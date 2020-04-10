package main

import "strings"

type Op interface{}

// HelpOp describes printing help.
type HelpOp struct{}

// ListOp describes listing contexts.
type ListOp struct{}

// CurrentOp prints the current context
type CurrentOp struct{}

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

// UnsetOp indicates intention to remove current-context preference.
type UnsetOp struct{}

// UnknownOp indicates an unsupported flag.
type UnknownOp struct{ Args []string }

// parseArgs looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func parseArgs(argv []string) Op {
	if len(argv) == 0 {
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
		if v == "--unset" || v == "-u" {
			return UnsetOp{}
		}

		if strings.HasPrefix(v, "-") && v != "-" {
			return UnknownOp{argv}
		}

		// TODO handle -d
		// TODO handle -u/--unset
		// TODO handle -c/--current
		return SwitchOp{Target: argv[0]}
	}

	// TODO handle too many arguments e.g. "kubectx a b c"
	return UnknownOp{}
}
