package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattn/go-isatty"
)

const (
	envFZFIgnore = "KUBECTX_IGNORE_FZF"
)

type InteractiveSwitchOp struct {
	SelfCmd string
}

func (op InteractiveSwitchOp) Run(_, stderr io.Writer) error {
	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", op.SelfCmd),
		fmt.Sprintf("%s=1", envForceColor))
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return errors.New("you did not choose any of the options")
	}
	name, err := switchContext(choice)
	if err != nil {
		return errors.Wrap(err, "failed to switch context")
	}
	printSuccess(stderr, "Switched to context %q.", name)
	return nil
}

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

// isInteractiveMode determines if we can do choosing with fzf.
func isInteractiveMode(stdout *os.File) bool {
	v := os.Getenv(envFZFIgnore)
	return v == "" && isTerminal(stdout) && fzfInstalled()
}
