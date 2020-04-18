package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (s SwitchOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(cmdutil.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	curNS, err := kc.NamespaceOfContext(ctx)
	if ctx == "" {
		return errors.New("failed to get current namespace")
	}

	f := NewNSFile(ctx)
	prev, err := f.Load()
	if err != nil {
		return errors.Wrap(err, "failed to load previous namespace from file")
	}

	toNS := s.Target
	if s.Target == "-" {
		if prev == "" {
			return errors.Errorf("No previous namespace found for current context (%s)", ctx)
		}
		toNS = prev
	}

	ok, err := namespaceExists(toNS)
	if err != nil {
		return errors.Wrap(err, "failed to query if namespace exists (is cluster accessible?)")
	}
	if !ok {
		return errors.Errorf("no namespace exists with name %q", toNS)
	}

	if err := kc.SetNamespace(ctx, toNS); err != nil {
		return errors.Wrapf(err, "failed to change to namespace %q", toNS)
	}
	if err := kc.Save(); err != nil {
		return errors.Wrap(err, "failed to save kubeconfig file")
	}
	if curNS != toNS {
		if err := f.Save(curNS); err != nil {
			return errors.Wrap(err, "failed to save the previous namespace to file")
		}
	}

	err = printer.Success(stderr, "Active namespace is %q", toNS)
	return err
}

func namespaceExists(ns string) (bool, error) {
	nses, err := queryNamespaces()
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
