package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// parseRenameSyntax parses A=B form into [A,B] and returns
// whether it is parsed correctly.
func parseRenameSyntax(v string) (string, string, bool) {
	s := strings.Split(v, "=")
	if len(s) != 2 {
		return "", "", false
	}
	new, old := s[0], s[1]
	if new == "" || old == "" {
		return "", "", false
	}
	return new, old, true
}

// rename changes the old (NAME or '.' for current-context)
// to the "new" value. If the old refers to the current-context,
// current-context preference is also updated.
func renameContexts(old, new string) error {
	f, rootNode, err := openKubeconfig()
	if err != nil {
		return nil
	}
	defer f.Close()

	cur := getCurrentContext(rootNode)
	if old == "." {
		old = cur
	}

	if !checkContextExists(rootNode, old) {
		return errors.Errorf("context %q not found, can't rename it", old)
	}

	if checkContextExists(rootNode, new) {
		printWarning("context %q exists, overwriting it.", new)
		if err := modifyDocToDeleteContext(rootNode, new); err != nil {
			return errors.Wrap(err, "failed to delete new context to overwrite it")
		}
	}

	if err := modifyContextName(rootNode, old, new); err != nil {
		return errors.Wrap(err, "failed to change context name")
	}

	if old == cur {
		if err := modifyCurrentContext(rootNode, new); err != nil {
			return errors.Wrap(err, "failed to set current-context to new name")
		}
	}

	// TODO the next two functions are always repeated.
	if err := resetFile(f); err != nil {
		return err
	}
	if err := saveKubeconfigRaw(f, rootNode); err != nil {
		return errors.Wrap(err, "failed to save modified kubeconfig")
	}
	return nil
}

func modifyContextName(rootNode *yaml.Node, old, new string) error {
	if rootNode.Kind != yaml.MappingNode {
		return errors.New("root doc is not a mapping node")
	}
	contexts := valueOf(rootNode, "contexts")
	if contexts == nil {
		return errors.New("\"contexts\" entry is nil")
	} else if contexts.Kind != yaml.SequenceNode {
		return errors.New("\"contexts\" is not a sequence node")
	}

	var changed bool
	for _, contextNode := range contexts.Content {
		nameNode := valueOf(contextNode, "name")
		if nameNode.Kind == yaml.ScalarNode && nameNode.Value == old {
			nameNode.Value = new
			changed = true
			break
		}
	}
	if !changed {
		return errors.New("no changes were made")
	}
	// TODO use printSuccess
	// TODO consider moving printing logic to main
	fmt.Fprintf(os.Stderr, "Context %q renamed to %q.\n", old, new)
	return nil
}
