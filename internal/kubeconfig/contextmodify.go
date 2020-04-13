package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) DeleteContextEntry(deleteName string) error {
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return errors.New("there are no contexts in kubeconfig")
	}
	if contexts.Kind != yaml.SequenceNode {
		return errors.New("'contexts' key is not a sequence")
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
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return errors.New("\"contexts\" entry is nil")
	} else if contexts.Kind != yaml.SequenceNode {
		return errors.New("\"contexts\" is not a sequence node")
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
