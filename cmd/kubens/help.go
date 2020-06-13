package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// HelpOp describes printing help.
type HelpOp struct{}

func (_ HelpOp) Run(stdout, _ io.Writer) error {
	return printUsage(stdout)
}

func printUsage(out io.Writer) error {
	help := `USAGE:
  %PROG%                    : list the namespaces in the current context
  %PROG% <NAME>             : change the active namespace of current context
  %PROG% -                  : switch to the previous namespace in this context
  %PROG% -c, --current      : show the current namespace
  %PROG% -h,--help          : show this message
  %PROG% -v,--version       : print the client version information
`
	// TODO this replace logic is duplicated between this and kubectx
	help = strings.ReplaceAll(help, "%PROG%", selfName())

	_, err := fmt.Fprintf(out, "%s\n", help)
	return errors.Wrap(err, "write error")
}

// selfName guesses how the user invoked the program.
func selfName() string {
	// TODO this method is duplicated between this and kubectx
	me := filepath.Base(os.Args[0])
	pluginPrefix := "kubectl-"
	if strings.HasPrefix(me, pluginPrefix) {
		return "kubectl " + strings.TrimPrefix(me, pluginPrefix)
	}
	return "kubens"
}
