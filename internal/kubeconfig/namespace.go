package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	defaultNamespace = "default"
)

func (k *Kubeconfig) contextNode(name string) (*yaml.Node, error) {
	contexts := valueOf(k.rootNode, "contexts")
	if contexts == nil {
		return nil, errors.New("\"contexts\" entry is nil")
	} else if contexts.Kind != yaml.SequenceNode {
		return nil, errors.New("\"contexts\" is not a sequence node")
	}

	for _, contextNode := range contexts.Content {
		nameNode := valueOf(contextNode, "name")
		if nameNode.Kind == yaml.ScalarNode && nameNode.Value == name {
			return contextNode, nil
		}
	}
	return nil, errors.Errorf("context with name %q not found", name)
}

func (k *Kubeconfig) NamespaceOfContext(contextName string) (string, error) {
	ctx, err := k.contextNode(contextName)
	if err != nil {
		return "", err
	}
	ns := valueOf(ctx, "namespace")
	if ns == nil || ns.Value == "" {
		return defaultNamespace, nil
	}
	return ns.Value, nil
}
