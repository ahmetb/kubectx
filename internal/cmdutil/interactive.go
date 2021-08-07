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

package cmdutil

import (
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"

	"github.com/ahmetb/kubectx/internal/env"
)

// isTerminal determines if given fd is a TTY.
func isTerminal(fd *os.File) bool {
	return isatty.IsTerminal(fd.Fd())
}

// pickerInstalled determines if picker(fzf or sk) is in PATH.
func pickerInstalled(p string) bool {
	v, _ := exec.LookPath(p)
	return v != ""
}

// IsInteractiveMode determines the picker and whether we can do interactive choosing
// with it.
func IsInteractiveMode(stdout *os.File) (string, bool) {
	p := fuzzyPicker()
	if p == "fzf" {
		v := os.Getenv(env.EnvFZFIgnore)
		return p, v == "" && isTerminal(stdout) && pickerInstalled(p)
	}
	// if picker is sk
	v := os.Getenv(env.EnvSKIgnore)
	return p, v == "" && isTerminal(stdout) && pickerInstalled(p)
}

func fuzzyPicker() string {
	p := os.Getenv(env.EnvPicker)
	if p == "sk" {
		return p
	}
	// for now it only supports fzf and sk.
	return "fzf"
}
