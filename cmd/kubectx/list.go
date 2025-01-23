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
	"errors"
	"fmt"
	"io"
	"io/fs"

	"facette.io/natsort"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return fmt.Errorf("kubeconfig error: %w", err)
	}

	ctxs := kc.ContextNames()
	natsort.Sort(ctxs)

	cur := kc.GetCurrentContext()
	for _, c := range ctxs {
		s := c
		if c == cur {
			s = printer.ActiveItemColor.Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}
