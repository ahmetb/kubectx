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

import "gopkg.in/yaml.v3"

const (
	defaultNamespace = "default"
)

func (k *Kubeconfig) NamespaceOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}
	ctxBody := valueOf(ctx, "context")
	if ctxBody == nil {
		return defaultNamespace, nil
	}
	ns := valueOf(ctxBody, "namespace")
	if ns == nil || ns.Value == "" {
		return defaultNamespace, nil
	}
	return ns.Value, nil
}

func (k *Kubeconfig) SetNamespace(ctxName string, ns string) error {
	ctxNode, err := k.contextNode(ctxName)
	if err != nil {
		return err
	}

	var ctxBodyNodeWasEmpty bool // actual namespace value is in contexts[index].context.namespace, but .context might not exist
	ctxBodyNode := valueOf(ctxNode, "context")
	if ctxBodyNode == nil {
		ctxBodyNodeWasEmpty = true
		ctxBodyNode = &yaml.Node{
			Kind: yaml.MappingNode,
		}
	}

	nsNode := valueOf(ctxBodyNode, "namespace")
	if nsNode != nil {
		nsNode.Value = ns
		return nil
	}

	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "namespace",
		Tag:   "!!str"}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: ns,
		Tag:   "!!str"}
	ctxBodyNode.Content = append(ctxBodyNode.Content, keyNode, valueNode)
	if ctxBodyNodeWasEmpty {
		ctxNode.Content = append(ctxNode.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "context",
			Tag:   "!!str",
		}, ctxBodyNode)
	}
	return nil
}
