package main

import (
	"fmt"
	"io"

	"facette.io/natsort"
	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
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
	natsort.Sort(ctxs)

	cur := kc.GetCurrentContext()
	for _, c := range ctxs {
		s := c
		if c == cur {
			s = printer.ActiveItemColor.Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}
