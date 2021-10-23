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
	"strings"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

type InteractiveSwitchOp struct {
	SelfCmd string
	Picker  string
}

// TODO(ahmetb) This method is heavily repetitive vs kubectx/fzf.go.
func (op InteractiveSwitchOp) Run(_, stderr io.Writer) error {
	// parse kubeconfig just to see if it can be loaded
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return errors.Wrap(err, "kubeconfig error")
	}
	defer kc.Close()
	out, err := cmdutil.InteractiveSearch(op.Picker, op.SelfCmd, stderr)
	if err != nil {
		return err
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return errors.New("you did not choose any of the options")
	}
	name, err := switchNamespace(kc, choice)
	if err != nil {
		return errors.Wrap(err, "failed to switch namespace")
	}
	printer.Success(stderr, "Active namespace is \"%s\".", printer.SuccessColor.Sprint(name))
	return nil
}
