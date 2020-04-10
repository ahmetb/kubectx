package main

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// switchContext switches to specified context name.
func switchContext(name string) (string, error) {
	stateFile, err := kubectxFilePath()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}

	f, kc, err := openKubeconfig()
	if err != nil {
		return "", err
	}
	defer f.Close()

	prev := getCurrentContext(kc)

	// TODO: add a check to ensure user can't switch to non-existing context.
	if !checkContextExists(kc, name)  {
		return "", errors.Errorf("no context exists with the name: %q", name)
	}

	if err := modifyCurrentContext(kc, name); err != nil {
		return "", err
	}

	if err := resetFile(f); err != nil {
		return "", err
	}

	if err := saveKubeconfigRaw(f, kc); err != nil {
		return "", errors.Wrap(err, "failed to save kubeconfig")
	}

	if prev != name {
		if err := writeLastContext(stateFile, prev); err != nil {
			return "", errors.Wrap(err, "failed to save previous context name")
		}
	}

	return name, nil
}

func checkContextExists(rootNode *yaml.Node, name string) bool {
	contexts := valueOf(rootNode, "contexts")
	if contexts == nil {
		return false
	}
	if contexts.Kind != yaml.SequenceNode {
		return false
	}

	var ctxNames []string
	for _, ctx := range contexts.Content {
		nameVal := valueOf(ctx, "name")
		if nameVal != nil {
			ctxNames = append(ctxNames, nameVal.Value)
		}
	}

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

// swapContext switches to previously switch context.
func swapContext() (string, error) {
	stateFile, err := kubectxFilePath()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}
	prev, err := readLastContext(stateFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to read previous context file")
	}
	if prev == "" {
		return "", errors.New("no previous context found")
	}
	return switchContext(prev)
}

// getCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func getCurrentContext(rootNode *yaml.Node) string {
	if rootNode.Kind != yaml.MappingNode {
		return ""
	}
	v := valueOf(rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
}

func modifyCurrentContext(rootNode *yaml.Node, name string) error {
	if rootNode.Kind != yaml.MappingNode {
		return errors.New("document is not a map")
	}

	// find current-context field => modify value (next children)
	for i, ch := range rootNode.Content {
		if i%2 == 0 && ch.Value == "current-context" {
			rootNode.Content[i+1].Value = name
			return nil
		}
	}

	// if current-context ==> create new field
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "current-context",
		Tag:   "!!str"}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: name,
		Tag:   "!!str"}
	rootNode.Content = append(rootNode.Content, keyNode, valueNode)
	return nil
}
