package printer

import (
	"os"

	"github.com/ahmetb/kubectx/internal/env"
)

// UseColors returns true if colors are force-enabled,
// false if colors are disabled, or nil for default behavior
// which is determined based on factors like if stdout is tty.
func UseColors() *bool {
	tr, fa := true, false
	if os.Getenv(env.EnvForceColor) != "" {
		return &tr
	} else if os.Getenv(env.EnvNoColor) != "" {
		return &fa
	}
	return nil
}
