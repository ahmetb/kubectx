package kubeconfig

import (
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) ContextNames() []string {
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return nil
	}
	if contexts.Kind != yaml.SequenceNode {
		return nil
	}

	var ctxNames []string
	for _, ctx := range contexts.Content {
		nameVal := valueOf(ctx, "name")
		if nameVal != nil {
			ctxNames = append(ctxNames, nameVal.Value)
		}
	}
	return ctxNames
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

func valueOf(mapNode *yaml.Node, key string) *yaml.Node {
	if mapNode.Kind != yaml.MappingNode {
		return nil
	}
	for i, ch := range mapNode.Content {
		if i%2 == 0 && ch.Kind == yaml.ScalarNode && ch.Value == key {
			return mapNode.Content[i+1]
		}
	}
	return nil
}
