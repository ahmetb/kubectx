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
	argv, tmpValue, tmpSet := stripTmpFlag(argv)

	if len(argv) == 0 {
		if cmdutil.IsInteractiveMode(os.Stdout) {
			return wrapTmpOp(InteractiveSwitchOp{SelfCmd: os.Args[0]}, tmpValue, tmpSet)
		}
		return wrapTmpOp(ListOp{}, tmpValue, tmpSet)
	}

	if argv[0] == "-d" {
		if len(argv) == 1 {
			if cmdutil.IsInteractiveMode(os.Stdout) {
				return wrapTmpOp(InteractiveDeleteOp{SelfCmd: os.Args[0]}, tmpValue, tmpSet)
			} else {
				return wrapTmpOp(UnsupportedOp{Err: fmt.Errorf("'-d' needs arguments")}, tmpValue, tmpSet)
			}
		}
		return wrapTmpOp(DeleteOp{Contexts: argv[1:]}, tmpValue, tmpSet)
	}

	if len(argv) == 1 {
		v := argv[0]
		if v == "--help" || v == "-h" {
			return wrapTmpOp(HelpOp{}, tmpValue, tmpSet)
		}
		if v == "--version" || v == "-V" {
			return wrapTmpOp(VersionOp{}, tmpValue, tmpSet)
		}
		if v == "--current" || v == "-c" {
			return wrapTmpOp(CurrentOp{}, tmpValue, tmpSet)
		}
		if v == "--unset" || v == "-u" {
			return wrapTmpOp(UnsetOp{}, tmpValue, tmpSet)
		}

		if new, old, ok := parseRenameSyntax(v); ok {
			return wrapTmpOp(RenameOp{New: new, Old: old}, tmpValue, tmpSet)
		}

		if strings.HasPrefix(v, "-") && v != "-" {
			return wrapTmpOp(UnsupportedOp{Err: fmt.Errorf("unsupported option '%s'", v)}, tmpValue, tmpSet)
		}
		return wrapTmpOp(SwitchOp{Target: argv[0]}, tmpValue, tmpSet)
	}
	return wrapTmpOp(UnsupportedOp{Err: fmt.Errorf("too many arguments")}, tmpValue, tmpSet)
}

func wrapTmpOp(op Op, value string, set bool) Op {
	if !set {
		return op
	}
	return TmpOp{Inner: op, Value: value}
}

func stripTmpFlag(argv []string) ([]string, string, bool) {
	var (
		out      []string
		tmpValue string
		tmpSet   bool
	)

	for _, v := range argv {
		if v == "--tmp" {
			tmpSet = true
			continue
		}
		if strings.HasPrefix(v, "--tmp=") {
			tmpSet = true
			tmpValue = strings.TrimPrefix(v, "--tmp=")
			continue
		}
		out = append(out, v)
	}
	return out, tmpValue, tmpSet
}
