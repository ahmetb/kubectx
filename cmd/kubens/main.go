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
	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/printer"
	"github.com/fatih/color"
)

type Op interface {
	Run(stdout, stderr io.Writer) error
}

var (
	name    = "kubens"
	version = "v0.0.0"
	date    = "0001-01-01T00:00:00Z"
	commit  = "0000000"
)

func ver() string {
	return fmt.Sprintf("%s %s (%s) %s", name, version, commit[:7], date)
}

func main() {
	cmdutil.PrintDeprecatedEnvWarnings(color.Error, os.Environ())

	// Support [--]version and -V
	if len(os.Args) > 1 {
		if "version" == strings.TrimLeft(os.Args[1], "-") || "-V" == os.Args[1] {
			fmt.Println(ver())
			os.Exit(0)
			return
		}
	}

	op := parseArgs(os.Args[1:])
	if err := op.Run(color.Output, color.Error); err != nil {
		printer.Error(color.Error, err.Error())

		if _, ok := os.LookupEnv(env.EnvDebug); ok {
			// print stack trace in verbose mode
			fmt.Fprintf(color.Error, "[DEBUG] error: %+v\n", err)
		}
		defer os.Exit(1)
	}
}
