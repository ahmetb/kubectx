package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

// UnsetOp indicates intention to remove current-context preference.
type UnsetOp struct{}

func (_ UnsetOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "failed to parse kubeconfig")
	}

	if err := kc.UnsetCurrentContext(); err != nil {
		return errors.Wrap(err, "error while modifying current-context")
	}
	if err := kc.Save(); err != nil {
		return errors.Wrap(err, "failed to save kubeconfig file after modification")
	}

	_, err := fmt.Fprintln(stderr, "Successfully unset the current context")
	return err
}
