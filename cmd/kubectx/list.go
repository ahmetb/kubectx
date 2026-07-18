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
	"fmt"
	"io"

	"facette.io/natsort"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	if err := checkIsolatedMode(); err != nil {
		return err
	}
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return fmt.Errorf("kubeconfig error: %w", err)
	}
	return formatContextList(kc, stdout)
}

// formatContextList writes the sorted, current-context-highlighted list of
// contexts from kc to w (one per line). The format mirrors what ListOp.Run
// prints so that the interactive fzf picker shows the same list the user would
// see from `kubectx | fzf`.
func formatContextList(kc *kubeconfig.Kubeconfig, w io.Writer) error {
	ctxs, err := kc.ContextNames()
	if err != nil {
		return fmt.Errorf("failed to get context names: %w", err)
	}
	natsort.Sort(ctxs)

	cur, err := kc.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("failed to get current context: %w", err)
	}
	for _, c := range ctxs {
		s := c
		if c == cur {
			s = printer.ActiveItemColor.Sprint(c)
		}
		fmt.Fprintf(w, "%s\n", s)
	}
	return nil
}
