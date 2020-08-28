package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (s SwitchOp) Run(_, stderr io.Writer) error {

	ok, err := namespaceExists(s.Target)
	if err != nil {
		return errors.Wrap(err, "failed to query if namespace exists (is cluster accessible?)")
	}
	if !ok {
		return errors.Errorf("no namespace exists with name %q", s.Target)
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	toNS, err := switchNamespace(kc, s.Target)
	if err != nil {
		return err
	}
	err = printer.Success(stderr, "Active namespace is %q", toNS)
	return err
}

func switchNamespace(kc *kubeconfig.Kubeconfig, ns string) (string, error) {

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return "", errors.New("current-context is not set")
	}

	curNS, err := kc.NamespaceOfContext(ctx)
	if curNS == "" {
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

func namespaceExists(ns string) (bool, error) {
	// for tests
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return ns == "ns1" || ns == "ns2" || ns == "-", nil
	}

	if ns == "-" {
		return true, nil
	}

	kubeCfgPath, err := kubeconfig.FindKubeconfigPath()
	if err != nil {
		return false, errors.Wrap(err, "cannot determine kubeconfig path")
	}

	clientset, err := newWritableKubernetesClientSet(kubeCfgPath)
	if err != nil {
		return false, errors.Wrap(err, "failed to initialize k8s REST client")
	}

	namespace, err := clientset.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if errors2.IsNotFound(err) {
		return false, nil
	}
	return namespace != nil, errors.Wrapf(err, "failed to query namespace %q from k8s API", ns)
}
