package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// switchContext switches to specified context name.
func switchContext(name string) (string, error) {
	stateFile, err := kubectxFilePath()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}

	cfgPath, err := kubeconfigPath()
	if err != nil {
		return "", errors.Wrap(err, "cannot determine kubeconfig path")
	}
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	kc, err := parseKubeconfigRaw(f)
	if err != nil {
		return "", errors.Wrap(err, "yaml parse error")
	}

	prev := getCurrentContext(kc)

	if err := modifyCurrentContext(kc, name); err != nil {
		return "", err
	}

	if  err := f.Truncate(0); err != nil {
		return "", errors.Wrap(err, "failed to truncate")
	}

	if _, err := f.Seek(0, 0); err != nil {
		return "", errors.Wrap(err, "failed to seek")
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


// getCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func getCurrentContext(rootNode *yaml.Node) string {
	if rootNode.Kind != yaml.MappingNode {
		return ""
	}
	for i, ch := range rootNode.Content {
		if i%2 == 0 && ch.Value == "current-context" {
			return rootNode.Content[i+1].Value
		}
	}
	return ""
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

func parseKubeconfigRaw(r io.Reader) (*yaml.Node, error) {
	var v yaml.Node
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return v.Content[0], nil
}

func saveKubeconfigRaw(w io.Writer, rootNode *yaml.Node) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(rootNode)
}
