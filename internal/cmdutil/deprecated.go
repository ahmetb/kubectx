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
	"io"
	"strings"

	"github.com/ahmetb/kubectx/internal/printer"
)

func PrintDeprecatedEnvWarnings(out io.Writer, vars []string) {
	for _, vv := range vars {
		parts := strings.SplitN(vv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]

		if key == `KUBECTX_CURRENT_FGCOLOR` || key == `KUBECTX_CURRENT_BGCOLOR` {
			printer.Warning(out, "%s environment variable is now deprecated", key)
		}
	}
}
