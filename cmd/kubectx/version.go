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
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// VersionOp describes printing help.
type VersionOp struct{}

func (_ VersionOp) Run(stdout, _ io.Writer) error {
	return printUsage(stdout)
}

func printUsage(out io.Writer) error {
	version := `verison: v0.10.0`
	_, err := fmt.Fprintf(out, "%s\n", version)
	return errors.Wrap(err, "write error")
}

// selfName guesses how the user invoked the program.
func selfName() string {
	me := filepath.Base(os.Args[0])
	pluginPrefix := "kubectl-"
	if strings.HasPrefix(me, pluginPrefix) {
		return "kubectl " + strings.TrimPrefix(me, pluginPrefix)
	}
	return "kubectx"
}
