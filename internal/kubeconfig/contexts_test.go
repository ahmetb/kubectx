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

	ctx, err := kc.ContextNames()
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"abc", "def", "ghi"}
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
	ctx, err := kc.ContextNames()
	if err != nil {
		t.Fatal(err)
	}
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
	_, err := kc.ContextNames()
	if err == nil {
		t.Fatal("expected error for non-array contexts entry")
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

	if exists, err := kc.ContextExists("c1"); err != nil || !exists {
		t.Fatal("c1 actually exists; reported false")
	}
	if exists, err := kc.ContextExists("c2"); err != nil || !exists {
		t.Fatal("c2 actually exists; reported false")
	}
	if exists, err := kc.ContextExists("c3"); err != nil {
		t.Fatal(err)
	} else if exists {
		t.Fatal("c3 does not exist; but reported true")
	}
}

func TestKubeconfig_ContextNames_MultiFile_Merged(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("a"), testutil.Ctx("b")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c"), testutil.Ctx("d")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	ctx, err := kc.ContextNames()
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"a", "b", "c", "d"}
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextNames_MultiFile_DuplicateFirstWins(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("a"), testutil.Ctx("shared")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("shared"), testutil.Ctx("b")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	ctx, err := kc.ContextNames()
	if err != nil {
		t.Fatal(err)
	}
	// "shared" should appear only once
	expected := []string{"a", "shared", "b"}
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextExists_MultiFile(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("c1")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if exists, err := kc.ContextExists("c1"); err != nil || !exists {
		t.Fatal("c1 should exist")
	}
	if exists, err := kc.ContextExists("c2"); err != nil || !exists {
		t.Fatal("c2 should exist (in second file)")
	}
	if exists, err := kc.ContextExists("c3"); err != nil {
		t.Fatal(err)
	} else if exists {
		t.Fatal("c3 should not exist")
	}
}
