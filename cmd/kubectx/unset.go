package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

// UnsetOp indicates intention to remove current-context preference.
type UnsetOp struct{}

func (_ UnsetOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()

	rootNode, err := kc.ParseRaw()
	if err != nil {
		return err
	}

	if err := modifyDocToUnsetContext(rootNode); err != nil {
		return errors.Wrap(err, "error while modifying current-context")
	}
	if err := kc.Save(); err != nil {
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
