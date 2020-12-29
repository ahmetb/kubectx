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

package kubeconfig

import (
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_GetCurrentContext(t *testing.T) {
	tl := WithMockKubeconfigLoader(`current-context: foo`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	v := kc.GetCurrentContext()

	expected := "foo"
	if v != expected {
		t.Fatalf("expected=\"%s\"; got=\"%s\"", expected, v)
	}
}

func TestKubeconfig_GetCurrentContext_missingField(t *testing.T) {
	tl := WithMockKubeconfigLoader(`abc: def`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	v := kc.GetCurrentContext()

	expected := ""
	if v != expected {
		t.Fatalf("expected=\"%s\"; got=\"%s\"", expected, v)
	}
}

func TestKubeconfig_UnsetCurrentContext(t *testing.T) {
	tl := WithMockKubeconfigLoader(testutil.KC().WithCurrentCtx("foo").ToYAML(t))
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.UnsetCurrentContext(); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	out := tl.Output()
	expected := testutil.KC().WithCurrentCtx("").ToYAML(t)
	if out != expected {
		t.Fatalf("expected=\"%s\"; got=\"%s\"", expected, out)
	}
}
