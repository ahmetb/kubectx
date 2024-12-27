// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type ListOp struct{}

func (op ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return fmt.Errorf("kubeconfig error, %w", err)
	}

	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	curNs, err := kc.NamespaceOfContext(ctx)
	if err != nil {
		return fmt.Errorf("cannot read current namespace, %w", err)
	}

	ns, err := queryNamespaces(kc)
	if err != nil {
		return fmt.Errorf("could not list namespaces (is the cluster accessible?), %w", err)
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

func queryNamespaces(kc *kubeconfig.Kubeconfig) ([]string, error) {
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return []string{"ns1", "ns2"}, nil
	}

	clientset, err := newKubernetesClientSet(kc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize k8s REST client, %w", err)
	}

	var out []string
	var next string
	for {
		list, err := clientset.CoreV1().Namespaces().List(
			context.Background(),
			metav1.ListOptions{
				Limit:    500,
				Continue: next,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces from k8s API, %w", err)
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

func newKubernetesClientSet(kc *kubeconfig.Kubeconfig) (*kubernetes.Clientset, error) {
	b, err := kc.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to convert in-memory kubeconfig to yaml, %w", err)
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config, %w", err)
	}
	return kubernetes.NewForConfig(cfg)
}
