package main

import (
	"fmt"
	"io"
)

func printHelp(out io.Writer) {
	help := `USAGE:
  kubectx                       : list the contexts
  kubectx <NAME>                : switch to context <NAME>
  kubectx -                     : switch to the previous context
  kubectx -c, --current         : show the current context name
  kubectx <NEW_NAME>=<NAME>     : rename context <NAME> to <NEW_NAME>
  kubectx <NEW_NAME>=.          : rename current-context to <NEW_NAME>
  kubectx -d <NAME> [<NAME...>] : delete context <NAME> ('.' for current-context)
                                  (this command won't delete the user/cluster entry
                                  that is used by the context)
  kubectx -u, --unset           : unset the current context

  kubectx -h,--help             : show this message`

	fmt.Fprintf(out, "%s\n", help)
}
