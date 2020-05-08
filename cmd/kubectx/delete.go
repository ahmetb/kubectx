package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
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
			return errors.Wrapf(err, "error deleting context %q", deletedName)
		}
		if wasActiveContext {
			printer.Warning(stderr, "You deleted the current context. Use \"%s\" to select a new context.",
				selfName())
		}

		printer.Success(stderr, `Deleted context %s.`, printer.SuccessColor.Sprint(deletedName))
	}
	return nil
}

// deleteContext deletes a context entry by NAME or current-context
// indicated by ".".
func deleteContext(name string) (deleteName string, wasActiveContext bool, err error) {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return deleteName, false, errors.Wrap(err, "kubeconfig error")
	}

	cur := kc.GetCurrentContext()
	// resolve "." to a real name
	if name == "." {
		if cur == "" {
			return deleteName, false, errors.New("can't use '.' as the no active context is set")
		}
		wasActiveContext = true
		name = cur
	}

	if !kc.ContextExists(name) {
		return name, false, errors.New("context does not exist")
	}

	if err := kc.DeleteContextEntry(name); err != nil {
		return name, false, errors.Wrap(err, "failed to modify yaml doc")
	}
	return name, wasActiveContext, errors.Wrap(kc.Save(), "failed to save modified kubeconfig file")
}
