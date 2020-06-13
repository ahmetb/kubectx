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
  %PROG%                       : list the contexts
  %PROG% <NAME>                : switch to context <NAME>
  %PROG% -                     : switch to the previous context
  %PROG% -c, --current         : show the current context name
  %PROG% <NEW_NAME>=<NAME>     : rename context <NAME> to <NEW_NAME>
  %PROG% <NEW_NAME>=.          : rename current-context to <NEW_NAME>
  %PROG% -u, --unset           : unset the current context  
  %PROG% -d <NAME> [<NAME...>] : delete context <NAME> ('.' for current-context)
  %SPAC%                         (this command won't delete the user/cluster entry
  %SPAC%                          referenced by the context entry)
  %PROG% -h,--help             : show this message
  %PROG% -v,--version          : print the client version information`
	help = strings.ReplaceAll(help, "%PROG%", selfName())
	help = strings.ReplaceAll(help, "%SPAC%", strings.Repeat(" ", len(selfName())))

	_, err := fmt.Fprintf(out, "%s\n", help)
	return errors.Wrap(err, "write error")
}

// selfName guesses how the user invoked the program.
func selfName() string {
	me := filepath.Base(os.Args[0])
	pluginPrefix := "kubectl-"
	if strings.HasPrefix(me, pluginPrefix) {
		return "kubectl " + strings.TrimPrefix(me, pluginPrefix)
	}
	return "kubectx"
}
