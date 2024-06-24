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
	"slices"
	"strings"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

// UnsupportedOp indicates an unsupported flag.
type UnsupportedOp struct{ Err error }

func (op UnsupportedOp) Run(_, _ io.Writer) error {
	return op.Err
}

// parseArgs looks at flags (excl. executable name, i.e. argv[0])
// and decides which operation should be taken.
func parseArgs(argv []string) Op {
	n := len(argv)

	if n == 0 {
		if cmdutil.IsInteractiveMode(os.Stdout) {
			return InteractiveSwitchOp{SelfCmd: os.Args[0]}
		}
		return ListOp{}
	}

	if n == 1 || n == 2 {
		v := argv[0]
		force := false

		if n == 2 {
			force = slices.Contains([]string{"-f", "--force"}, argv[1])
		}

		if v == "--help" || v == "-h" {
			return HelpOp{}
		}
		if v == "--version" || v == "-V" {
			return VersionOp{}
		}
		if v == "--current" || v == "-c" {
			return CurrentOp{}
		}
		if strings.HasPrefix(v, "-") && v != "-" {
			return UnsupportedOp{Err: fmt.Errorf("unsupported option '%s'", v)}
		}
		return SwitchOp{Target: argv[0], Force: force}
	}
	return UnsupportedOp{Err: fmt.Errorf("too many arguments")}
}
