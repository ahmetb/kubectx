package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

// DeleteOp indicates intention to delete contexts.
type DeleteOp struct {
	Contexts []string // NAME or '.' to indicate current-context.
}

// deleteContexts deletes context entries one by one.
func (op DeleteOp) Run(_, stderr io.Writer) error {
	for _, ctx := range op.Contexts {
		// TODO inefficency here. we open/write/close the same file many times.
		deletedName, wasActiveContext, err := deleteContext(ctx)
		if err != nil {
			return errors.Wrapf(err, "error deleting context %q", ctx)
		}
		if wasActiveContext {
			// TODO we don't always run as kubectx (sometimes "kubectl ctx")
			printWarning(stderr, "You deleted the current context. use \"kubectx\" to select a different one.")
		}

		printSuccess(stderr, "deleted context %q", deletedName)
	}
	return nil
}

// deleteContext deletes a context entry by NAME or current-context
// indicated by ".".
func deleteContext(name string) (deleteName string, wasActiveContext bool, err error) {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return "", false, errors.Wrap(err, "kubeconfig error")
	}

	cur := kc.GetCurrentContext()

	// resolve "." to a real name
	if name == "." {
		wasActiveContext = true
		name = cur
	}

	if !kc.ContextExists(name) {
		return "", false, errors.New("context does not exist")
	}

	if err := kc.DeleteContextEntry(name); err != nil {
		return "", false, errors.Wrap(err, "failed to modify yaml doc")
	}
	return name, wasActiveContext, errors.Wrap(kc.Save(), "failed to save modified kubeconfig file")
}
