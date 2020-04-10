package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func printCurrentContext(w io.Writer) error {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return errors.Wrap(err, "failed to determine kubeconfig path")
	}

	cfg, err := parseKubeconfig(cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig file")
	}

	v := cfg.CurrentContext
	if v == "" {
		return errors.New("current-context is not set")
	}
	_, err = fmt.Fprintln(w, v)
	return err
}
