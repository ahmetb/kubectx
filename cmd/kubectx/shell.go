package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ShellOp indicates intention to start a scoped sub-shell for a context.
type ShellOp struct {
	Target string
}

func (op ShellOp) Run(_, stderr io.Writer) error {
	if err := checkIsolatedMode(); err != nil {
		return err
	}

	kubectlPath, err := resolveKubectl()
	if err != nil {
		return err
	}

	// Verify context exists and get current context for exit message
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}
	if !kc.ContextExists(op.Target) {
		return fmt.Errorf("no context exists with the name: \"%s\"", op.Target)
	}
	previousCtx := kc.GetCurrentContext()

	// Extract minimal kubeconfig using kubectl
	data, err := extractMinimalKubeconfig(kubectlPath, op.Target)
	if err != nil {
		return errors.Wrap(err, "failed to extract kubeconfig for context")
	}

	// Write to temp file
	tmpFile, err := os.CreateTemp("", "kubectx-shell-*.yaml")
	if err != nil {
		return errors.Wrap(err, "failed to create temp kubeconfig file")
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return errors.Wrap(err, "failed to write temp kubeconfig")
	}
	tmpFile.Close()

	// Print entry message
	badgeColor := color.New(color.BgRed, color.FgWhite, color.Bold)
	printer.EnableOrDisableColor(badgeColor)
	fmt.Fprintf(stderr, "%s kubectl context is %s in this shell — type 'exit' to leave.\n",
		badgeColor.Sprint("[ISOLATED SHELL]"), printer.WarningColor.Sprint(op.Target))

	// Detect and start shell
	shellBin := detectShell()
	cmd := exec.Command(shellBin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+tmpPath,
		env.EnvIsolatedShell+"=1",
	)

	_ = cmd.Run()

	// Print exit message
	fmt.Fprintf(stderr, "%s kubectl context is now %s.\n",
		badgeColor.Sprint("[ISOLATED SHELL EXITED]"), printer.WarningColor.Sprint(previousCtx))

	return nil
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
