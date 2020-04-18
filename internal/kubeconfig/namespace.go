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
