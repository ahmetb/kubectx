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

package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

var (
	tr, fa = true, false
)

func Test_useColors_forceColors(t *testing.T) {
	defer testutil.WithEnvVar("_KUBECTX_FORCE_COLOR", "1")()
	defer testutil.WithEnvVar("NO_COLOR", "1")()

	if v := useColors(); !cmp.Equal(v, &tr) {
		t.Fatalf("expected useColors() = true; got = %v", v)
	}
}

func Test_useColors_disableColors(t *testing.T) {
	defer testutil.WithEnvVar("NO_COLOR", "1")()

	if v := useColors(); !cmp.Equal(v, &fa) {
		t.Fatalf("expected useColors() = false; got = %v", v)
	}
}

func Test_useColors_default(t *testing.T) {
	defer testutil.WithEnvVar("NO_COLOR", "")()
	defer testutil.WithEnvVar("_KUBECTX_FORCE_COLOR", "")()

	if v := useColors(); v != nil {
		t.Fatalf("expected useColors() = nil; got=%v", *v)
	}
}
