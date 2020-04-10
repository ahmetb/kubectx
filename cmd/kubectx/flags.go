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

// DeleteOp indicates intention to delete contexts.
type DeleteOp struct {
	Contexts []string // NAME or '.' to indicate current-context.
}

// RenameOp indicates intention to rename contexts.
type RenameOp struct {
	New string // NAME of New context
	Old string // NAME of Old context (or '.' for current-context)
}

// UnknownOp indicates an unsupported flag.
type UnknownOp struct{ Args []string }

// parseArgs looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func parseArgs(argv []string) Op {
	if len(argv) == 0 {
		return ListOp{}
	}

	if argv[0] == "-d" {
		ctxs := argv[1:]
		return DeleteOp{ctxs}
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

		new, old, ok := parseRenameSyntax(v) // a=b a=.
		if ok {
			return RenameOp{new, old}
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
