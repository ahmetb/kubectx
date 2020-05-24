package main

import (
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// RenameOp indicates intention to rename contexts.
type RenameOp struct {
	New string // NAME of New context
	Old string // NAME of Old context (or '.' for current-context)
}

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
func (op RenameOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}

	cur := kc.GetCurrentContext()
	if op.Old == "." {
		op.Old = cur
	}

	if !kc.ContextExists(op.Old) {
		return errors.Errorf("context \"%s\" not found, can't rename it", op.Old)
	}

	if kc.ContextExists(op.New) {
		printer.Warning(stderr, "context \"%s\" exists, overwriting it.", op.New)
		if err := kc.DeleteContextEntry(op.New); err != nil {
			return errors.Wrap(err, "failed to delete new context to overwrite it")
		}
	}

	if err := kc.ModifyContextName(op.Old, op.New); err != nil {
		return errors.Wrap(err, "failed to change context name")
	}
	if op.Old == cur {
		if err := kc.ModifyCurrentContext(op.New); err != nil {
			return errors.Wrap(err, "failed to set current-context to new name")
		}
	}
	if err := kc.Save(); err != nil {
		return errors.Wrap(err, "failed to save modified kubeconfig")
	}
	printer.Success(stderr, "Context %s renamed to %s.",
		printer.SuccessColor.Sprint(op.Old),
		printer.SuccessColor.Sprint(op.New))
	return nil
}
