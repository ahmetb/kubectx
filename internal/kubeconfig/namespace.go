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

package kubeconfig

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	defaultNamespace = "default"
)

func (k *Kubeconfig) NamespaceOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}
	namespace, err := ctx.Pipe(yaml.Lookup("context", "namespace"))
	if namespace == nil || err != nil {
		return defaultNamespace, err
	}
	return yaml.GetValue(namespace), nil
}

func (k *Kubeconfig) SetNamespace(ctxName string, ns string) error {
	ctx, err := k.contextNode(ctxName)
	if err != nil {
		return err
	}
	if err := ctx.PipeE(
		yaml.LookupCreate(yaml.MappingNode, "context"),
		yaml.SetField("namespace", yaml.NewStringRNode(ns)),
	); err != nil {
		return err
	}
	return nil
}
