package main

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

// SwitchOp indicates intention to switch contexts.
type SwitchOp struct {
	Target string // '-' for back and forth, or NAME
}

func (op SwitchOp) Run(stdout, stderr io.Writer) error {
	var newCtx string
	var err error
	if op.Target == "-" {
		newCtx, err = swapContext()
	} else {
		newCtx, err = switchContext(op.Target)
	}
	if err != nil {
		return errors.Wrap(err, "failed to switch context")
	}
	printSuccess(stderr, "Switched to context %q.", newCtx)
	return nil
}

// switchContext switches to specified context name.
func switchContext(name string) (string, error) {
	prevCtxFile, err := kubectxPrevCtxFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}

	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()

	rootNode, err := kc.ParseRaw()
	if err != nil {
		return "", err
	}

	prev := kubeconfig.GetCurrentContext(rootNode)
	if !checkContextExists(rootNode, name) {
		return "", errors.Errorf("no context exists with the name: %q", name)
	}
	if err := modifyCurrentContext(rootNode, name); err != nil {
		return "", err
	}
	if err := kc.Save(); err != nil {
		return "", errors.Wrap(err, "failed to save kubeconfig")
	}

	if prev != name {
		if err := writeLastContext(prevCtxFile, prev); err != nil {
			return "", errors.Wrap(err, "failed to save previous context name")
		}
	}
	return name, nil
}


// swapContext switches to previously switch context.
func swapContext() (string, error) {
	prevCtxFile, err := kubectxPrevCtxFile()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine state file")
	}
	prev, err := readLastContext(prevCtxFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to read previous context file")
	}
	if prev == "" {
		return "", errors.New("no previous context found")
	}
	return switchContext(prev)
}


func checkContextExists(rootNode *yaml.Node, name string) bool {
	ctxNames := kubeconfig.ContextNames(rootNode)
	for _, v := range ctxNames {
		if v == name {
			return true
		}
	}
	return false
}

// TODO delete
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

	// if current-context ==> create New field
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
