package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// shellSession holds the configuration for spawning an isolated sub-shell.
type shellSession struct {
	target   string
	extraEnv []string // additional env vars beyond KUBECONFIG + KUBECTX_ISOLATED_SHELL

	printEntry func(stderr io.Writer, ctxName string)
	printExit  func(stderr io.Writer, prevCtx string)

	// transformKubeconfig optionally transforms the minified kubeconfig bytes
	// before writing them to the shell's temp file. The returned cleanup func
	// is called after the shell exits (e.g. to shut down a proxy).
	// If nil, the kubeconfig is used as-is.
	transformKubeconfig func(data []byte) (newData []byte, cleanup func(), err error)
}

func (s *shellSession) run(stderr io.Writer) error {
	if err := checkIsolatedMode(); err != nil {
		return err
	}

	kubectlPath, err := resolveKubectl()
	if err != nil {
		return err
	}

	// Verify context exists and get current context for exit message.
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return fmt.Errorf("kubeconfig error: %w", err)
	}
	exists, err := kc.ContextExists(s.target)
	if err != nil {
		return fmt.Errorf("failed to check context: %w", err)
	}
	if !exists {
		return fmt.Errorf("no context exists with the name: %q", s.target)
	}
	previousCtx, err := kc.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to get current context: %w", err)
	}

	// Extract minimal kubeconfig for the target context.
	data, err := extractMinimalKubeconfig(kubectlPath, s.target)
	if err != nil {
		return fmt.Errorf("failed to extract kubeconfig for context: %w", err)
	}

	// Optionally transform the kubeconfig (e.g. rewrite for readonly proxy).
	var cleanup func()
	if s.transformKubeconfig != nil {
		data, cleanup, err = s.transformKubeconfig(data)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}
	}

	// Write kubeconfig to temp file.
	tmpFile, err := os.CreateTemp("", "kubectx-shell-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp kubeconfig file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp kubeconfig: %w", err)
	}
	tmpFile.Close()

	// Print entry message.
	s.printEntry(stderr, s.target)

	// Detect and start shell.
	shellBin := detectShell()
	cmd := exec.Command(shellBin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		append([]string{
			"KUBECONFIG=" + tmpPath,
			env.EnvIsolatedShell + "=1",
		}, s.extraEnv...)...,
	)

	_ = cmd.Run()

	// Print exit message.
	s.printExit(stderr, previousCtx)

	return nil
}

// fzfPickContext launches fzf for interactive context selection.
func fzfPickContext(selfCmd string, stderr io.Writer) (string, error) {
	if err := checkIsolatedMode(); err != nil {
		return "", err
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return "", nil
		}
		return "", fmt.Errorf("kubeconfig error: %w", err)
	}

	ctxNames, err := kc.ContextNames()
	if err != nil {
		return "", fmt.Errorf("failed to get context names: %w", err)
	}
	if len(ctxNames) == 0 {
		return "", errors.New("no contexts found in the kubeconfig file")
	}

	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", selfCmd),
		fmt.Sprintf("%s=1", env.EnvForceColor))
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return "", err
		}
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return "", errors.New("you did not choose any of the options")
	}
	return choice, nil
}
