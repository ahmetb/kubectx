package kubeconfig

import (
	"gopkg.in/yaml.v3"
)

func deleteNamedChildNode(node *yaml.Node, childName string) {
	i := -1
	for j, node := range node.Content {
		nameNode := valueOf(node, "name")
		if nameNode != nil && nameNode.Kind == yaml.ScalarNode && nameNode.Value == childName {
			i = j
			break
		}
	}

	if i >= 0 {
		copy(node.Content[i:], node.Content[i+1:])
		node.Content[len(node.Content)-1] = nil
		node.Content = node.Content[:len(node.Content)-1]
	}
}
