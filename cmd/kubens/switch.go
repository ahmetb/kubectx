package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type SwitchOp struct {
	Target  string // '-' for back and forth, or NAME
	IsForce bool
}

func (s SwitchOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(cmdutil.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	toNS, err := switchNamespace(kc, s.Target, s.IsForce)
	if err != nil {
		return err
	}
	err = printer.Success(stderr, "Active namespace is %q", toNS)
	return err
}

func switchNamespace(kc *kubeconfig.Kubeconfig, ns string, isForce bool) (string, error) {
	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return "", errors.New("current-context is not set")
	}
	curNS, err := kc.NamespaceOfContext(ctx)
	if ctx == "" {
		return "", errors.New("failed to get current namespace")
	}

	f := NewNSFile(ctx)
	prev, err := f.Load()
	if err != nil {
		return "", errors.Wrap(err, "failed to load previous namespace from file")
	}

	if ns == "-" {
		if prev == "" {
			return "", errors.Errorf("No previous namespace found for current context (%s)", ctx)
		}
		ns = prev
	}

	if !isForce {
		ok, err := namespaceExists(kc, ns)
		if err != nil {
			return "", errors.Wrap(err, "failed to query if namespace exists (is cluster accessible?)")
		}
		if !ok {
			return "", errors.Errorf("no namespace exists with name %q", ns)
		}
	}

	if err := kc.SetNamespace(ctx, ns); err != nil {
		return "", errors.Wrapf(err, "failed to change to namespace %q", ns)
	}
	if err := kc.Save(); err != nil {
		return "", errors.Wrap(err, "failed to save kubeconfig file")
	}
	if curNS != ns {
		if err := f.Save(curNS); err != nil {
			return "", errors.Wrap(err, "failed to save the previous namespace to file")
		}
	}
	return ns, nil
}

func namespaceExists(kc *kubeconfig.Kubeconfig, ns string) (bool, error) {
	nses, err := queryNamespaces(kc)
	if err != nil {
		return false, err
	}
	for _, v := range nses {
		if v == ns {
			return true, nil
		}
	}
	return false, nil
}
