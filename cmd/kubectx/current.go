package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

// CurrentOp prints the current context
type CurrentOp struct{}

func (_op CurrentOp) Run(stdout, _ io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()
	_, err := kc.ParseRaw()
	if err != nil {
		return err
	}

	v := kc.GetCurrentContext()
	if v == "" {
		return errors.New("current-context is not set")
	}
	_, err = fmt.Fprintln(stdout, v)
	return err
}
