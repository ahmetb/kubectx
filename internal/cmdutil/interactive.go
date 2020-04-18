package env

import (
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"

	"github.com/ahmetb/kubectx/internal/env"
)

// isTerminal determines if given fd is a TTY.
func isTerminal(fd *os.File) bool {
	return isatty.IsTerminal(fd.Fd())
}

// fzfInstalled determines if fzf(1) is in PATH.
func fzfInstalled() bool {
	v, _ := exec.LookPath("fzf")
	if v != "" {
		return true
	}
	return false
}

// IsInteractiveMode determines if we can do choosing with fzf.
func IsInteractiveMode(stdout *os.File) bool {
	v := os.Getenv(env.EnvFZFIgnore)
	return v == "" && isTerminal(stdout) && fzfInstalled()
}
