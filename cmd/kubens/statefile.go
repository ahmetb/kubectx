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
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ahmetb/kubectx/core/cmdutil"
)

var defaultDir = filepath.Join(cmdutil.HomeDir(), ".kube", "kubens")

type NSFile struct {
	dir string
	ctx string
}

func NewNSFile(ctx string) NSFile { return NSFile{dir: defaultDir, ctx: ctx} }

func (f NSFile) path() string {
	fn := f.ctx
	if isWindows() {
		// bug 230: eks clusters contain ':' in ctx name, not a valid file name for win32
		fn = strings.ReplaceAll(fn, ":", "__")
	}
	return filepath.Join(f.dir, fn)
}

// Load reads the previous namespace setting, or returns empty if not exists.
func (f NSFile) Load() (string, error) {
	b, err := ioutil.ReadFile(f.path())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

// Save stores the previous namespace information in the file.
func (f NSFile) Save(value string) error {
	d := filepath.Dir(f.path())
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(f.path(), []byte(value), 0644)
}

// isWindows determines if the process is running on windows OS.
func isWindows() bool {
	if os.Getenv("_FORCE_GOOS") == "windows" { // for testing
		return true
	}
	return runtime.GOOS == "windows"
}
