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

	if n == 1 {
		v := argv[0]
		switch v {
		case "--help", "-h":
			return HelpOp{}
		case "--version", "-V":
			return VersionOp{}
		case "--current", "-c":
			return CurrentOp{}
		default:
			return getSwitchOp(v, false)
		}
	} else if n == 2 {
		// {namespace} -f|--force
		name := argv[0]
		force := slices.Contains([]string{"-f", "--force"}, argv[1])

		if !force {
			if !slices.Contains([]string{"-f", "--force"}, argv[0]) {
				return UnsupportedOp{Err: fmt.Errorf("unsupported arguments %q", argv)}
			}

			// -f|--force {namespace}
			force = true
			name = argv[1]
		}

		return getSwitchOp(name, force)
	}

	return UnsupportedOp{Err: fmt.Errorf("too many arguments")}
}

func getSwitchOp(v string, force bool) Op {
	if strings.HasPrefix(v, "-") && v != "-" {
		return UnsupportedOp{Err: fmt.Errorf("unsupported option %q", v)}
	}
	return SwitchOp{Target: v, Force: force}
}
