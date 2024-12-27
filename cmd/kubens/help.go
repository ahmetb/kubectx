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
)

// HelpOp describes printing help.
type HelpOp struct{}

func (_ HelpOp) Run(stdout, _ io.Writer) error {
	return printUsage(stdout)
}

func printUsage(out io.Writer) error {
	help := `USAGE:
  %PROG%                    : list the namespaces in the current context
  %PROG% <NAME>             : change the active namespace of current context
  %PROG% <NAME> --force/-f  : force change the active namespace of current context (even if it doesn't exist)
  %PROG% -                  : switch to the previous namespace in this context
  %PROG% -c, --current      : show the current namespace
  %PROG% -h,--help          : show this message
  %PROG% -V,--version       : show version`

	// TODO this replace logic is duplicated between this and kubectx
	help = strings.ReplaceAll(help, "%PROG%", selfName())

	_, err := fmt.Fprintf(out, "%s\n", help)
	if err != nil {
		return fmt.Errorf("write error, %w", err)
	}
	return nil
}

// selfName guesses how the user invoked the program.
func selfName() string {
	// TODO this method is duplicated between this and kubectx
	me := filepath.Base(os.Args[0])
	pluginPrefix := "kubectl-"
	if strings.HasPrefix(me, pluginPrefix) {
		return "kubectl " + strings.TrimPrefix(me, pluginPrefix)
	}
	return "kubens"
}
