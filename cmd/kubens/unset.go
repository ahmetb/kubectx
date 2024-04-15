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

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
	"github.com/pkg/errors"
)

// UnsetOp indicates intention to remove current namespace preference.
type UnsetOp struct{}

func (u UnsetOp) Run(_, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return errors.Wrap(err, "kubeconfig error")
	}
	ctx := kc.GetCurrentContext()
	if ctx == "" {
		return errors.New("current-context is not set")
	}
	err := kc.UnsetNamespace(ctx)
	if err != nil {
		return errors.Wrap(err, "could not unset namespace for the current-context")
	}
	if err := kc.Save(); err != nil {
		return errors.Wrap(err, "failed to save kubeconfig file")
	}

	err = printer.Success(stderr, "Active namespace is unset")
	return err
}
