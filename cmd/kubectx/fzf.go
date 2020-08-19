package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type InteractiveSwitchOp struct {
	SelfCmd string
}

func (op InteractiveSwitchOp) Run(_, stderr io.Writer) error {
	// parse kubeconfig just to see if it can be loaded
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)

	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return errors.Wrap(err, "kubeconfig error")
	}

	ctxs := kc.ContextNames()
	if ctxs == nil {
		err := printer.Warning(stderr, "No kubectl context found")
		return errors.Wrap(err, "kubeconfig error")
	}
	kc.Close()

	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", op.SelfCmd),
		fmt.Sprintf("%s=1", env.EnvForceColor))
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
	printer.Success(stderr, "Switched to context %s.", printer.SuccessColor.Sprint(name))
	return nil
}
