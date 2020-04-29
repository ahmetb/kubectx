package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type ListOp struct{}

func (op ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(cmdutil.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	curNs, err := kc.NamespaceOfContext(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot read current namespace")
	}

	ns, err := queryNamespaces()
	if err != nil {
		return errors.Wrap(err, "could not list namespaces (is the cluster accessible?)")
	}


	for _, c := range ns {
		s := c
		if c == curNs {
			s = printer.ActiveItemColor.Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}

func findKubectl() (string, error) {
	if v := os.Getenv("KUBECTL"); v != "" {
		return v, nil
	}
	v, err := exec.LookPath("kubectl")
	return v, errors.Wrap(err, "kubectl not found, needed for kubens")
}

func queryNamespaces() ([]string, error) {
	kubectl ,err := findKubectl()
	if err != nil {
		return nil ,err
	}

	// TODO add a log message to user if kubectl is taking >1s

	var b bytes.Buffer
	cmd := exec.Command(kubectl, "get", "namespaces", `-o=jsonpath={range .items[*].metadata.name}{@}{"\n"}{end}`)
	cmd.Env = os.Environ()
	cmd.Stdout, cmd.Stderr = &b, &b
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "failed to query namespaces: %v", b.String())
	}
	return strings.Split(strings.TrimSpace(b.String()), "\n"), nil
}
