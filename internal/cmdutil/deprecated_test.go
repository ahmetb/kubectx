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
	"bytes"
	"strings"
	"testing"
)

func TestPrintDeprecatedEnvWarnings_noDeprecatedVars(t *testing.T) {
	var out bytes.Buffer
	PrintDeprecatedEnvWarnings(&out, []string{
		"A=B",
		"PATH=/foo:/bar:/bin",
	})
	if v := out.String(); len(v) > 0 {
		t.Fatalf("something written to buf: %v", v)
	}
}

func TestPrintDeprecatedEnvWarnings_bgColors(t *testing.T) {
	var out bytes.Buffer

	PrintDeprecatedEnvWarnings(&out, []string{
		"KUBECTX_CURRENT_FGCOLOR=1",
		"KUBECTX_CURRENT_BGCOLOR=2",
	})
	v := out.String()
	if !strings.Contains(v, "KUBECTX_CURRENT_FGCOLOR") {
		t.Fatalf("output doesn't contain 'KUBECTX_CURRENT_FGCOLOR': \"%s\"", v)
	}
	if !strings.Contains(v, "KUBECTX_CURRENT_BGCOLOR") {
		t.Fatalf("output doesn't contain 'KUBECTX_CURRENT_BGCOLOR': \"%s\"", v)
	}
}
