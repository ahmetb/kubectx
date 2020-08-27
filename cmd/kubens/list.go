package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type ListOp struct{}

func (op ListOp) Run(stdout, _ io.Writer) error {

	allNamespaces, err := queryNamespaces()
	if err != nil {
		return errors.Wrap(err, "could not list namespaces (is the cluster accessible?)")
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
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

	for _, c := range allNamespaces {
		s := c
		if c == curNs {
			s = printer.ActiveItemColor.Sprint(c)
		}
		_, _ = fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}

func queryNamespaces() ([]string, error) {

	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return []string{"ns1", "ns2"}, nil
	}

	kubeCfgPath, err := kubeconfig.FindKubeconfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}

	clientset, err := newWritableKubernetesClientSet(kubeCfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize k8s REST client")
	}

	var out []string
	var next string
	for {
		list, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{
			Limit:    500,
			Continue: next,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list namespaces from k8s API")
		}
		next = list.Continue
		for _, it := range list.Items {
			out = append(out, it.Name)
		}
		if next == "" {
			break
		}
	}
	return out, nil
}

func newWritableKubernetesClientSet(kubeCfgPath string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeCfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize REST client config")
	}

	return kubernetes.NewForConfig(config)
}
