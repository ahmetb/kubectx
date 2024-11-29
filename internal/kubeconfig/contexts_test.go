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

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_ContextNames(t *testing.T) {
	tl := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("abc"),
			testutil.Ctx("def"),
			testutil.Ctx("ghi")).Set("field1", map[string]string{"bar": "zoo"}).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	ctx := kc.ContextNames(nil)
	expected := []string{"abc", "def", "ghi"}
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextNames_Filtered(t *testing.T) {
	tl := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("abc").Ns("ns1"),
			testutil.Ctx("def"),
			testutil.Ctx("ghi").Ns("ns2"),
		).Set("field1", map[string]string{"bar": "zoo"}).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	ctx := kc.ContextNames(map[string]string{"namespace": "ns2"})
	expected := []string{"ghi"}
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextNames_noContextsEntry(t *testing.T) {
	tl := WithMockKubeconfigLoader(`a: b`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	ctx := kc.ContextNames(nil)
	var expected []string = nil
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextNames_nonArrayContextsEntry(t *testing.T) {
	tl := WithMockKubeconfigLoader(`contexts: "hello"`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	ctx := kc.ContextNames(nil)
	var expected []string = nil
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_CheckContextExists(t *testing.T) {
	tl := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2")).ToYAML(t))

	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if !kc.ContextExists("c1") {
		t.Fatal("c1 actually exists; reported false")
	}
	if !kc.ContextExists("c2") {
		t.Fatal("c2 actually exists; reported false")
	}
	if kc.ContextExists("c3") {
		t.Fatal("c3 does not exist; but reported true")
	}
}
