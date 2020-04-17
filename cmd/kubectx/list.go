package main

import (
	"fmt"
	"io"
	"os"

	"facette.io/natsort"
	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if isENOENT(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return errors.Wrap(err, "kubeconfig error")
	}

	ctxs := kc.ContextNames()
	natsort.Sort(ctxs)

	// TODO support KUBECTX_CURRENT_FGCOLOR
	// TODO support KUBECTX_CURRENT_BGCOLOR

	currentColor := color.New(color.FgGreen, color.Bold)
	if v := printer.UseColors(); v != nil && *v {
		currentColor.EnableColor()
	} else if v != nil && !*v {
		currentColor.DisableColor()
	}

	cur := kc.GetCurrentContext()
	for _, c := range ctxs {
		s := c
		if c == cur {
			s = currentColor.Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}

// isENOENT determines if the underlying error is os.IsNotExist. Right now
// errors from github.com/pkg/errors doesn't work with os.IsNotExist.
func isENOENT(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if os.IsNotExist(e) {
			return true
		}
	}
	return false
}
