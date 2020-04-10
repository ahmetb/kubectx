package main

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func unsetContext() error {
	f, rootNode, err := openKubeconfig()
	if err != nil {
		return err
	}
	defer f.Close()

	if err := modifyDocToUnsetContext(rootNode); err != nil {
		return errors.Wrap(err, "error while modifying current-context")
	}
	if err := resetFile(f); err != nil {
		return err
	}
	if err := saveKubeconfigRaw(f, rootNode); err != nil {
		return errors.Wrap(err, "failed to save kubeconfig file after modification")
	}
	return nil
}

func modifyDocToUnsetContext(rootNode *yaml.Node) error {
	if rootNode.Kind != yaml.MappingNode {
		return errors.New("kubeconfig file is not a map document")
	}
	curCtxValNode := valueOf(rootNode, "current-context")
	curCtxValNode.Value = ""
	return nil
}
