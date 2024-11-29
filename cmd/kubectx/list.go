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
	"strings"

	"facette.io/natsort"
	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
	"github.com/ahmetb/kubectx/internal/printer"
)

// ListOp describes listing contexts.
type ListOp struct {
	Filters map[string]string
}

// parseFilterSyntax parses multiple A=B form into a map[A]=B and returns
// whether it is parsed correctly.
func parseFilterSyntax(v []string) (map[string]string, bool) {
	m := make(map[string]string)
	for _, vv := range v {
		s := strings.Split(vv, "=")
		if len(s) != 2 {
			return nil, false
		}
		key, value := s[0], s[1]
		if key == "" || value == "" {
			return nil, false
		}
		m[key] = value
	}
	return m, true
}

func (op ListOp) Run(stdout, stderr io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		if cmdutil.IsNotFoundErr(err) {
			printer.Warning(stderr, "kubeconfig file not found")
			return nil
		}
		return errors.Wrap(err, "kubeconfig error")
	}

	ctxs := kc.ContextNames(op.Filters)
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
