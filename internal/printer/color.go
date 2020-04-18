package printer

import (
	"os"

	"github.com/fatih/color"

	"github.com/ahmetb/kubectx/internal/env"
)

// useColors returns true if colors are force-enabled,
// false if colors are disabled, or nil for default behavior
// which is determined based on factors like if stdout is tty.
func useColors() *bool {
	tr, fa := true, false
	if os.Getenv(env.EnvForceColor) != "" {
		return &tr
	} else if os.Getenv(env.EnvNoColor) != "" {
		return &fa
	}
	return nil
}

// EnableOrDisableColor determines if color should be force-enabled or force-disabled
// or left untouched based on environment configuration.
func EnableOrDisableColor(c *color.Color) {
	if v := useColors(); v != nil && *v {
		c.EnableColor()
	} else if v != nil && !*v {
		c.DisableColor()
	}
}
