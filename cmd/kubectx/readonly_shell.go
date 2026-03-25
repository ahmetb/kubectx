package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
	"github.com/ahmetb/kubectx/internal/proxy"
)

// InteractiveReadonlyShellOp launches fzf to pick a context, then starts a readonly shell.
type InteractiveReadonlyShellOp struct {
	SelfCmd string
}

// ReadonlyShellOp starts a read-only sub-shell for a context.
type ReadonlyShellOp struct {
	Target string
}

func (op InteractiveReadonlyShellOp) Run(_, stderr io.Writer) error {
	if err := checkIsolatedMode(); err != nil {
		return err
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return fmt.Errorf("kubeconfig error: %w", err)
	}

	ctxNames, err := kc.ContextNames()
	if err != nil {
		return fmt.Errorf("failed to get context names: %w", err)
	}
	if len(ctxNames) == 0 {
		return errors.New("no contexts found in the kubeconfig file")
	}

	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", op.SelfCmd),
		fmt.Sprintf("%s=1", env.EnvForceColor))
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return err
		}
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return errors.New("you did not choose any of the options")
	}
	return ReadonlyShellOp{Target: choice}.Run(nil, stderr)
}

func (op ReadonlyShellOp) Run(_, stderr io.Writer) error {
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
	exists, err := kc.ContextExists(op.Target)
	if err != nil {
		return fmt.Errorf("failed to check context: %w", err)
	}
	if !exists {
		return fmt.Errorf("no context exists with the name: %q", op.Target)
	}
	previousCtx, err := kc.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to get current context: %w", err)
	}

	// Extract minimal kubeconfig for the target context.
	data, err := extractMinimalKubeconfig(kubectlPath, op.Target)
	if err != nil {
		return fmt.Errorf("failed to extract kubeconfig for context: %w", err)
	}

	// Write original minified kubeconfig to temp file (used by proxy for TLS/auth).
	origFile, err := os.CreateTemp("", "kubectx-readonly-orig-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp kubeconfig file: %w", err)
	}
	origPath := origFile.Name()
	defer os.Remove(origPath)

	if _, err := origFile.Write(data); err != nil {
		origFile.Close()
		return fmt.Errorf("failed to write temp kubeconfig: %w", err)
	}
	origFile.Close()

	// Start the readonly proxy.
	p, err := proxy.Start(proxy.Config{
		KubeconfigPath: origPath,
		ContextName:    op.Target,
	})
	if err != nil {
		return fmt.Errorf("failed to start readonly proxy: %w", err)
	}
	defer p.Shutdown(context.Background())

	// Rewrite kubeconfig to point to the proxy.
	rewritten, err := proxy.RewriteKubeconfig(data, p.Addr())
	if err != nil {
		return fmt.Errorf("failed to rewrite kubeconfig: %w", err)
	}

	// Write rewritten kubeconfig to a second temp file for the shell.
	shellFile, err := os.CreateTemp("", "kubectx-readonly-shell-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp kubeconfig file: %w", err)
	}
	shellPath := shellFile.Name()
	defer os.Remove(shellPath)

	if _, err := shellFile.Write(rewritten); err != nil {
		shellFile.Close()
		return fmt.Errorf("failed to write rewritten kubeconfig: %w", err)
	}
	shellFile.Close()

	// Give the proxy a moment to be ready.
	time.Sleep(10 * time.Millisecond)

	// Print entry message.
	badgeColor := color.New(color.BgYellow, color.FgBlack, color.Bold)
	printer.EnableOrDisableColor(badgeColor)
	fmt.Fprintf(stderr, "%s kubectl context is %s in READ-ONLY mode — type 'exit' to leave.\n",
		badgeColor.Sprint("[READONLY SHELL]"), printer.WarningColor.Sprint(op.Target))

	// Detect and start shell.
	shellBin := detectShell()
	cmd := exec.Command(shellBin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"KUBECONFIG="+shellPath,
		env.EnvIsolatedShell+"=1",
		env.EnvReadonlyShell+"=1",
	)

	_ = cmd.Run()

	// Print exit message.
	fmt.Fprintf(stderr, "%s kubectl context is now %s.\n",
		badgeColor.Sprint("[READONLY SHELL EXITED]"), printer.WarningColor.Sprint(previousCtx))

	return nil
}
