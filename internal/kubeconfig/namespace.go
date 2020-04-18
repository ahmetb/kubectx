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
	ns := valueOf(ctx, "namespace")
	if ns == nil || ns.Value == "" {
		return defaultNamespace, nil
	}
	return ns.Value, nil
}

func (k *Kubeconfig) SetNamespace(ctxName string, ns string) error {
	ctx, err := k.contextNode(ctxName)
	if err != nil {
		return err
	}
	nsNode := valueOf(ctx, "namespace")
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
	ctx.Content = append(ctx.Content, keyNode, valueNode)
	return nil
}
