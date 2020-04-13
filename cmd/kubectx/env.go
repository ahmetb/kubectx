package main

import "os"

const (
	// EnvFZFIgnore describes the environment variable to set to disable
	// interactive context selection when fzf is installed.
	EnvFZFIgnore = "KUBECTX_IGNORE_FZF"

	// EnvForceColor describes the environment variable to disable color usage
	// when printing current context in a list.
	EnvNoColor = `NO_COLOR`

	// EnvForceColor describes the "internal" environment variable to force
	// color usage to show current context in a list.
	EnvForceColor = `_KUBECTX_FORCE_COLOR`

	// EnvDebug describes the internal environment variable for more verbose logging.
	EnvDebug = `DEBUG`
)

func useColors() bool {
	if os.Getenv(EnvForceColor) != "" {
		return true
	} else if os.Getenv(EnvNoColor) != "" {
		return false
	}
	return true
}
