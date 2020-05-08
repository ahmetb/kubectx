package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

// CurrentOp prints the current context
type CurrentOp struct{}

func (_op CurrentOp) Run(stdout, _ io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	v := kc.GetCurrentContext()
	if v == "" {
		return errors.New("current-context is not set")
	}
	_, err := fmt.Fprintln(stdout, v)
	return errors.Wrap(err, "write error")
}
