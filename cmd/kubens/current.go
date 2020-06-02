package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

type CurrentOp struct{}

func (c CurrentOp) Run(stdout, _ io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	ns, err := kc.NamespaceOfContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to read namespace of %q", ctx)
	}
	_, err = fmt.Fprintln(stdout, ns)
	return errors.Wrap(err, "write error")
}
