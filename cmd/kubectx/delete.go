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

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// DeleteOp indicates intention to delete contexts.
type DeleteOp struct {
	Contexts []string // NAME or '.' to indicate current-context.
	Cascade  bool     // Whether to delete (orphaned-only) users and clusters referenced in the contexts.
}

// deleteContexts deletes context entries one by one.
func (op DeleteOp) Run(_, stderr io.Writer) error {
	for _, ctx := range op.Contexts {
		// TODO inefficency here. we open/write/close the same file many times.
		deletedName, wasActiveContext, err := deleteContext(ctx, op.Cascade)
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
// The cascade flag determines whether to also delete the user and/or cluster entries referenced in the context,
// if they became orphaned by this deletion (i.e., not referenced by any other contexts).
func deleteContext(name string, cascade bool) (deleteName string, wasActiveContext bool, err error) {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return deleteName, false, errors.Wrap(err, "kubeconfig error")
	}

	cur := kc.GetCurrentContext()
	// resolve "." to a real name
	if name == "." {
		if cur == "" {
			return deleteName, false, errors.New("can't use '.' as no active context is set")
		}
		name = cur
	}

	wasActiveContext = name == cur

	if !kc.ContextExists(name) {
		return name, false, errors.New("context does not exist")
	}

	if cascade {
		err = deleteContextUser(name, kc)
		if err != nil {
			return name, wasActiveContext, errors.Wrap(err, "failed to delete user for deleted context")
		}

		err = deleteContextCluster(name, kc)
		if err != nil {
			return name, wasActiveContext, errors.Wrap(err, "failed to delete cluster for deleted context")
		}
	}

	if err := kc.DeleteContextEntry(name); err != nil {
		return name, false, errors.Wrap(err, "failed to modify yaml doc")
	}

	return name, wasActiveContext, errors.Wrap(kc.Save(), "failed to save modified kubeconfig file")
}

func deleteContextUser(contextName string, kc *kubeconfig.Kubeconfig) error {
	userName, err := kc.UserOfContext(contextName)
	if err != nil {
		return errors.Wrap(err, "user not set for context")
	}

	refCount, err := kc.CountUserReferences(userName)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve reference count for user entry")
	}

	if refCount == 1 {
		if err := kc.DeleteUserEntry(userName); err != nil {
			return errors.Wrap(err, "failed to modify yaml doc")
		}
	}

	return nil
}

func deleteContextCluster(contextName string, kc *kubeconfig.Kubeconfig) error {
	clusterName, err := kc.ClusterOfContext(contextName)
	if err != nil {
		return errors.Wrap(err, "cluster not set for context")
	}

	refCount, err := kc.CountClusterReferences(clusterName)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve reference count for cluster entry")
	}

	if refCount == 1 {
		if err := kc.DeleteClusterEntry(clusterName); err != nil {
			return errors.Wrap(err, "failed to modify yaml doc")
		}
	}

	return nil
}
