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
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type ListOp struct {
	cache *kubeNsCache
}

type stringSet map[string]struct{}

func (op *ListOp) LoadCache() error {
	if op.cache != nil {
		return nil
	}

	cache, err := newKubensCache()

	if err != nil {
		cache := make(kubeNsCache)
		op.cache = &cache

		return fmt.Errorf("failed to load cache: %w", err)
	}

	op.cache = &cache

	return nil
}

func (op ListOp) Run(stdout, stderr io.Writer) error {
	if err := op.LoadCache(); err != nil {
		if _, ok := os.LookupEnv(env.EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(color.Error, "[DEBUG] failed to load cache: %+v\n", err)
		}
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

	nsSet := make(stringSet)

	curNs, err := kc.NamespaceOfContext(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot read current namespace")
	}

	var namespaces []string

	namespaces, exists := op.cache.Get(ctx)
	if exists {
		// pass the cached list for now
		for _, ns := range namespaces {
			nsSet[ns] = struct{}{}

			if ns == curNs {
				ns = printer.ActiveItemColor.Sprint(ns)
			}

			fmt.Fprintf(stdout, "%s\n", ns)
		}
	}

	// now we start querying for the new list

	namespaces, err = queryNamespaces(kc)
	if err != nil {
		return errors.Wrap(err, "could not list namespaces (is the cluster accessible?)")
	}

	for _, ns := range namespaces {
		_, exists := nsSet[ns]
		if exists {
			continue
		}
		// didn't have this before, append to stdin
		fmt.Fprintf(stdout, "%s\n", ns)
	}

	op.cache.Save(ctx, namespaces)

	return nil
}

func queryNamespaces(kc *kubeconfig.Kubeconfig) ([]string, error) {
	if os.Getenv("_MOCK_NAMESPACES") != "" {
		return []string{"ns1", "ns2"}, nil
	}

	clientset, err := newKubernetesClientSet(kc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize k8s REST client")
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

func newKubernetesClientSet(kc *kubeconfig.Kubeconfig) (*kubernetes.Clientset, error) {
	b, err := kc.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert in-memory kubeconfig to yaml")
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize config")
	}
	return kubernetes.NewForConfig(cfg)
}
