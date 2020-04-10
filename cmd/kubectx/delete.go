package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// deleteContexts deletes context entries one by one.
func deleteContexts(w io.Writer, ctxs []string) error {
	for _, ctx := range ctxs {
		// TODO inefficency here. we open/write/close the same file many times.
		deletedName, wasActiveContext, err := deleteContext(ctx)
		if err != nil {
			return errors.Wrapf(err, "error deleting context %q", ctx)
		}
		if wasActiveContext {
			// TODO we don't always run as kubectx (sometimes "kubectl ctx")
			printWarning("You deleted the current context. use \"kubectx\" to select a different one.")
		}
		fmt.Fprintf(w, "deleted context %q\n", deletedName) // TODO write with printSuccess (i.e. green)
	}
	return nil
}

// deleteContext deletes a context entry by NAME or current-context
// indicated by ".".
func deleteContext(name string) (deleteName string, wasActiveContext bool, err error) {
	f, rootNode, err := openKubeconfig()
	if err != nil {
		return "", false, err
	}
	defer f.Close()

	cur := getCurrentContext(rootNode)

	// resolve "." to a real name
	if name == "." {
		wasActiveContext = true
		name = cur
	}

	if !checkContextExists(rootNode, name) {
		return "", false, errors.New("context does not exist")
	}

	if err := modifyDocToDeleteContext(rootNode, name); err != nil {
		return "", false, errors.Wrap(err, "failed to modify yaml doc")
	}

	if err := resetFile(f); err != nil {
		return "", false, err
	}
	return name, wasActiveContext, errors.Wrap(saveKubeconfigRaw(f, rootNode), "failed to save kubeconfig file")
}

func modifyDocToDeleteContext(rootNode *yaml.Node, deleteName string) error {
	if rootNode.Kind != yaml.MappingNode {
		return errors.New("root node was not a mapping node")
	}
	contexts := valueOf(rootNode, "contexts")
	if contexts == nil {
		return errors.New("there are no contexts in kubeconfig")
	}
	if contexts.Kind != yaml.SequenceNode {
		return errors.New("'contexts' key is not a sequence")
	}

	i := -1
	for j, ctxNode := range contexts.Content {
		nameNode := valueOf(ctxNode, "name")
		if nameNode != nil && nameNode.Kind == yaml.ScalarNode && nameNode.Value == deleteName {
			i = j
			break
		}
	}
	if i >= 0 {
		copy(contexts.Content[i:], contexts.Content[i+1:])
		contexts.Content[len(contexts.Content)-1] = nil
		contexts.Content = contexts.Content[:len(contexts.Content)-1]
	}
	return nil
}
