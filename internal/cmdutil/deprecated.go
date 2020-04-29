package cmdutil

import (
	"io"
	"strings"

	"github.com/ahmetb/kubectx/internal/printer"
)

func PrintDeprecatedEnvWarnings(out io.Writer, vars []string) {
	for _, vv := range vars {
		parts := strings.SplitN(vv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]

		if key == `KUBECTX_CURRENT_FGCOLOR` || key == `KUBECTX_CURRENT_BGCOLOR` {
			printer.Warning(out, "%s environment variable is now deprecated", key)
		}
	}
}
