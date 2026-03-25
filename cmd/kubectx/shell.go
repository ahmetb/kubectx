package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"

	"github.com/ahmetb/kubectx/internal/printer"
)

// InteractiveShellOp launches fzf to pick a context, then starts an isolated shell.
type InteractiveShellOp struct {
	SelfCmd string
}

// ShellOp indicates intention to start a scoped sub-shell for a context.
type ShellOp struct {
	Target string
}

func (op InteractiveShellOp) Run(_, stderr io.Writer) error {
	choice, err := fzfPickContext(op.SelfCmd, stderr)
	if err != nil || choice == "" {
		return err
	}
	return ShellOp{Target: choice}.Run(nil, stderr)
}

func (op ShellOp) Run(_, stderr io.Writer) error {
	badgeColor := color.New(color.BgRed, color.FgWhite, color.Bold)
	printer.EnableOrDisableColor(badgeColor)

	s := &shellSession{
		target: op.Target,
		printEntry: func(w io.Writer, ctxName string) {
			fmt.Fprintf(w, "%s kubectl context is %s in this shell — type 'exit' to leave.\n",
				badgeColor.Sprint("[ISOLATED SHELL]"), printer.WarningColor.Sprint(ctxName))
		},
		printExit: func(w io.Writer, prevCtx string) {
			fmt.Fprintf(w, "%s kubectl context is now %s.\n",
				badgeColor.Sprint("[ISOLATED SHELL EXITED]"), printer.WarningColor.Sprint(prevCtx))
		},
	}
	return s.run(stderr)
}

func resolveKubectl() (string, error) {
	if v := os.Getenv("KUBECTL"); v != "" {
		return v, nil
	}
	path, err := exec.LookPath("kubectl")
	if err != nil {
		return "", fmt.Errorf("kubectl is required for --shell but was not found in PATH")
	}
	return path, nil
}

func extractMinimalKubeconfig(kubectlPath, contextName string) ([]byte, error) {
	cmd := exec.Command(kubectlPath, "config", "view", "--minify", "--flatten",
		"--context", contextName)
	cmd.Env = os.Environ()
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("kubectl config view failed: %w", err)
	}
	return data, nil
}

func detectShell() string {
	if runtime.GOOS == "windows" {
		// cmd.exe always sets the PROMPT env var, so if it is present
		// we can reliably assume we are running inside cmd.exe.
		if os.Getenv("PROMPT") != "" {
			return "cmd.exe"
		}
		// Otherwise assume PowerShell. PSModulePath is always set on
		// Windows regardless of the shell, so it cannot be used as a
		// discriminator; however the absence of PROMPT is a strong
		// enough signal that we are in a PowerShell session.
		if pwsh, err := exec.LookPath("pwsh"); err == nil {
			return pwsh
		}
		if powershell, err := exec.LookPath("powershell"); err == nil {
			return powershell
		}
		return "cmd.exe"
	}
	if v := os.Getenv("SHELL"); v != "" {
		return v
	}
	return "/bin/sh"
}
