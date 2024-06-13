// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/ahmetb/kubectx/core/printer"
)

// DeleteOp indicates intention to delete contexts.
type DeleteOp struct {
	Contexts []string // NAME or '.' to indicate current-context.
}

// deleteContexts deletes context entries one by one.
func (op DeleteOp) Run(_, stderr io.Writer) error {
	for _, ctx := range op.Contexts {
		// TODO inefficiency here. we open/write/close the same file many times.
		deletedName, wasActiveContext, err := deleteContext(ctx)
		if err != nil {
			return errors.Wrapf(err, "error deleting context \"%s\"", deletedName)
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
