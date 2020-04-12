package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// UnsetOp indicates intention to remove current-context preference.
type UnsetOp struct{}

func (_ UnsetOp) Run(_, stderr io.Writer) error {
	f, rootNode, err := openKubeconfig()
	if err != nil {
		return err
	}
	defer f.Close()

	if err := modifyDocToUnsetContext(rootNode); err != nil {
		return errors.Wrap(err, "error while modifying current-context")
	}
	if err := saveKubeconfigRaw(f, rootNode); err != nil {
		return errors.Wrap(err, "failed to save kubeconfig file after modification")
	}

	_, err = fmt.Fprintln(stderr, "Successfully unset the current context")
	return err
}

func modifyDocToUnsetContext(rootNode *yaml.Node) error {
	if rootNode.Kind != yaml.MappingNode {
		return errors.New("kubeconfig file is not a map document")
	}
	curCtxValNode := valueOf(rootNode, "current-context")
	curCtxValNode.Value = ""
	return nil
}
