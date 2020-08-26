package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (op SwitchOp) Run(_, stderr io.Writer) error {
	var newCtx string
	var err error
	if op.Target == "-" {
		newCtx, err = swapContext()
	} else {
		newCtx, err = switchContext(op.Target)
	}
	if err != nil {
		return errors.Wrap(err, "failed to switch context")
	}
	err = printer.Success(stderr, "Switched to context \"%s\".", printer.SuccessColor.Sprint(newCtx))
	return errors.Wrap(err, "print error")
}

// switchContext switches to specified context name.
func switchContext(name string) (string, error) {
	prevCtxFile, err := kubectxPrevCtxFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return "", errors.Wrap(err, "kubeconfig error")
	}

	prev := kc.GetCurrentContext()
	if !kc.ContextExists(name) {
		return "", errors.Errorf("no context exists with the name: \"%s\"", name)
	}
	if err := kc.ModifyCurrentContext(name); err != nil {
		return "", err
	}
	if err := kc.Save(); err != nil {
		return "", errors.Wrap(err, "failed to save kubeconfig")
	}

	if prev != name {
		if err := writeLastContext(prevCtxFile, prev); err != nil {
			return "", errors.Wrap(err, "failed to save previous context name")
		}
	}
	return name, nil
}

// swapContext switches to previously switch context.
func swapContext() (string, error) {
	prevCtxFile, err := kubectxPrevCtxFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}
	prev, err := readLastContext(prevCtxFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to read previous context file")
	}
	if prev == "" {
		return "", errors.New("no previous context found")
	}
	return switchContext(prev)
}
