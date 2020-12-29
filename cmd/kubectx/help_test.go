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
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer
	if err := (&HelpOp{}).Run(&buf, &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "USAGE:") {
		t.Errorf("help string doesn't contain USAGE: ; output=\"%s\"", out)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("does not end with New line; output=\"%s\"", out)
	}
}
