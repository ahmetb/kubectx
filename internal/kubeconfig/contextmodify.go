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
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) DeleteContextEntry(deleteName string) error {
	contexts, err := k.contextsNode()
	if err != nil {
		return err
	}

	i := -1
	for j, ctxNode := range contexts.Content {
		nameNode := valueOf(ctxNode, "name")
		if nameNode != nil && nameNode.Kind == yaml.ScalarNode && nameNode.Value == deleteName {
			i = j
			break
		}
	}
	if i >= 0 {
		copy(contexts.Content[i:], contexts.Content[i+1:])
		contexts.Content[len(contexts.Content)-1] = nil
		contexts.Content = contexts.Content[:len(contexts.Content)-1]
	}
	return nil
}

func (k *Kubeconfig) ModifyCurrentContext(name string) error {
	currentCtxNode := valueOf(k.rootNode, "current-context")
	if currentCtxNode != nil {
		currentCtxNode.Value = name
		return nil
	}

	// if current-context field doesn't exist, create new field
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "current-context",
		Tag:   "!!str"}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: name,
		Tag:   "!!str"}
	k.rootNode.Content = append(k.rootNode.Content, keyNode, valueNode)
	return nil
}

func (k *Kubeconfig) ModifyContextName(old, new string) error {
	contexts, err := k.contextsNode()
	if err != nil {
		return err
	}

	var changed bool
	for _, contextNode := range contexts.Content {
		nameNode := valueOf(contextNode, "name")
		if nameNode.Kind == yaml.ScalarNode && nameNode.Value == old {
			nameNode.Value = new
			changed = true
			break
		}
	}
	if !changed {
		return errors.New("no changes were made")
	}
	return nil
}
