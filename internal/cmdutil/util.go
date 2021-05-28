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

	"github.com/pkg/errors"
)

func HomeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}
	return home
}

// IsNotFoundErr determines if the underlying error is os.IsNotExist. Right now
// errors from github.com/pkg/errors doesn't work with os.IsNotExist.
func IsNotFoundErr(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if os.IsNotExist(e) {
			return true
		}
	}
	return false
}
