package kubeconfig

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// GetCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func  (k *Kubeconfig) GetCurrentContext() string {
	if k.rootNode.Kind != yaml.MappingNode {
		return ""
	}
	v := valueOf(k.rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
}

func (k *Kubeconfig) UnsetCurrentContext() error {
	if k.rootNode.Kind != yaml.MappingNode {
		return errors.New("kubeconfig file is not a map document")
	}
	curCtxValNode := valueOf(k.rootNode, "current-context")
	curCtxValNode.Value = ""
	return nil
}
