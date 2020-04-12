package main

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

type Op interface {
	Run(stdout, stderr io.Writer) error
}

// UnsupportedOp indicates an unsupported flag.
type UnsupportedOp struct{ Err error }

func (op UnsupportedOp) Run(_, _ io.Writer) error {
	return op.Err
}

// parseArgs looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func parseArgs(argv []string) Op {
	if len(argv) == 0 {
		return ListOp{}
	}

	if argv[0] == "-d" {
		return DeleteOp{Contexts: argv[1:]}
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

		if new, old, ok := parseRenameSyntax(v); ok {
			return RenameOp{New: new, Old: old}
		}

		if strings.HasPrefix(v, "-") && v != "-" {
			return UnsupportedOp{Err: errors.Errorf("unsupported option %s", v)}
		}
		return SwitchOp{Target: argv[0]}
	}

	// TODO handle too many arguments e.g. "kubectx a b c"
	return UnsupportedOp{Err: errors.New("too many arguments")}
}
