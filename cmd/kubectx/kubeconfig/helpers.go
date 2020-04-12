package kubeconfig

import "gopkg.in/yaml.v3"

func ContextNames(rootNode *yaml.Node) []string {
	contexts := valueOf(rootNode, "contexts")
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

// GetCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func GetCurrentContext(rootNode *yaml.Node) string {
	if rootNode.Kind != yaml.MappingNode {
		return ""
	}
	v := valueOf(rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
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
