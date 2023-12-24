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
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (k *Kubeconfig) contextsNode() (*yaml.RNode, error) {
	contexts, err := k.config.Pipe(yaml.Get("contexts"))
	if err != nil {
		return nil, err
	}
	if contexts == nil {
		return nil, errors.New("\"contexts\" entry is nil")
	} else if contexts.YNode().Kind != yaml.SequenceNode {
		return nil, errors.New("\"contexts\" is not a sequence node")
	}
	return contexts, nil
}

func (k *Kubeconfig) contextNode(name string) (*yaml.RNode, error) {
	context, err := k.config.Pipe(yaml.Lookup("contexts", "[name="+name+"]"))
	if err != nil {
		return nil, err
	}
	if context == nil {
		return nil, errors.Errorf("context with name \"%s\" not found", name)
	}
	return context, nil
}

func (k *Kubeconfig) ContextNames() []string {
	contexts, err := k.config.Pipe(yaml.Get("contexts"))
	if err != nil {
		return nil
	}
	names, err := contexts.ElementValues("name")
	if err != nil {
		return nil
	}
	return names
}

func (k *Kubeconfig) ContextExists(name string) bool {
	ctxNames := k.ContextNames()
	for _, v := range ctxNames {
		if v == name {
			return true
		}
	}
	return false
}
