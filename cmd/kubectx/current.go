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

	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

// CurrentOp prints the current context
type CurrentOp struct{}

func (_op CurrentOp) Run(stdout, _ io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		return fmt.Errorf("kubeconfig error, %w", err)
	}

	v := kc.GetCurrentContext()
	if v == "" {
		return errors.New("current-context is not set")
	}
	_, err := fmt.Fprintln(stdout, v)
	if err != nil {
		return fmt.Errorf("write error, %w", err)
	}
	return nil
}
