package main

import (
	"io"

	"github.com/pkg/errors"
)

type context struct {
	Name string `yaml:"name"`
}

type kubeconfig struct {
	APIVersion     string    `yaml:"apiVersion"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []context `yaml:"contexts"`
}

func printListContexts(out io.Writer) error {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return errors.Wrap(err, "failed to determine kubeconfig path")
	}

	cfg, err := parseKubeconfig(cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig file")
	}
	_ = cfg
	
	// print each context
	//  - natural sort
	//  - highlight current context

	return nil
}
