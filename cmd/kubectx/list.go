package main

import (
	"fmt"
	"io"

	"facette.io/natsort"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type context struct {
	Name string `yaml:"name"`
}

type kubeconfigContents struct {
	APIVersion     string    `yaml:"apiVersion"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []context `yaml:"contexts"`
}

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	// TODO extract printing and sorting into a function that's testable

	cfgPath, err := kubeconfigPath()
	if err != nil {
		return errors.Wrap(err, "failed to determine kubeconfig path")
	}

	cfg, err := parseKubeconfig(cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig file")
	}

	ctxs := make([]string, 0, len(cfg.Contexts))
	for _, c := range cfg.Contexts {
		ctxs = append(ctxs, c.Name)
	}
	natsort.Sort(ctxs)

	// TODO support KUBECTX_CURRENT_FGCOLOR
	// TODO support KUBECTX_CURRENT_BGCOLOR
	for _, c := range ctxs {
		s := c
		if c == cfg.CurrentContext {
			s = color.New(color.FgGreen, color.Bold).Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}
