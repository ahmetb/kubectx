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

	ns, err := queryNamespaces(kc)
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

func getKubernetesClientForConfig(kc *kubeconfig.Kubeconfig) (*kubernetes.Clientset, error) {
	b, err := kc.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert in-memory kubeconfig to yaml")
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize config")
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize k8s REST client")
	}
	return clientset, nil
}

func queryNamespaces(kc *kubeconfig.Kubeconfig) ([]string, error) {
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return []string{"ns1", "ns2"}, nil
	}

	clientset, err := getKubernetesClientForConfig(kc)
	if err != nil {
		return nil, err
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

func namespaceExists(kc *kubeconfig.Kubeconfig, NamespaceName string) (bool, error) {
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		if NamespaceName == "ns1" || NamespaceName == "ns2" {
			return true, nil
		}
		return  false, nil
	}
	clientset, err := getKubernetesClientForConfig(kc)
	if err != nil {
		return false, err
	}
	_, err = clientset.CoreV1().Namespaces().Get(NamespaceName, metav1.GetOptions{})
	if err != nil {
		return false, errors.Wrap(err, "failed to get namespace from k8s API")
	}
	return true, nil
}
