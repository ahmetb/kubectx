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

func TestKubeconfig_NamespaceOfContext_ctxNotFound(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(testutil.Ctx("c1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	_, err := kc.NamespaceOfContext("c2")
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestKubeconfig_NamespaceOfContext(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2").Ns("c2n1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	v1, err := kc.NamespaceOfContext("c1")
	if err != nil {
		t.Fatal("expected err")
	}
	if expected := `default`; v1 != expected {
		t.Fatalf("c1: expected=\"%s\" got=\"%s\"", expected, v1)
	}

	v2, err := kc.NamespaceOfContext("c2")
	if err != nil {
		t.Fatal("expected err")
	}
	if expected := `c2n1`; v2 != expected {
		t.Fatalf("c2: expected=\"%s\" got=\"%s\"", expected, v2)
	}
}

func TestKubeconfig_SetNamespace(t *testing.T) {
	l := WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2").Ns("c2n1")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(l)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if err := kc.SetNamespace("c3", "foo"); err == nil {
		t.Fatalf("expected error for non-existing ctx")
	}

	if err := kc.SetNamespace("c1", "c1n1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.SetNamespace("c2", "c2n2"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithCtxs(
		testutil.Ctx("c1").Ns("c1n1"),
		testutil.Ctx("c2").Ns("c2n2")).ToYAML(t)
	if diff := cmp.Diff(l.Output(), expected); diff != "" {
		t.Fatal(diff)
	}
}

func TestKubeconfig_NamespaceOfContext_MultiFile(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("c1").Ns("ns1")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c2").Ns("ns2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	v, err := kc.NamespaceOfContext("c2")
	if err != nil {
		t.Fatal(err)
	}
	if v != "ns2" {
		t.Fatalf("expected=\"ns2\" got=\"%s\"", v)
	}
}

func TestKubeconfig_SetNamespace_MultiFile(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("c1")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	// Set namespace for c2 which is in second file
	if err := kc.SetNamespace("c2", "my-ns"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	// First file should be unchanged
	out0 := tl.OutputOf(0)
	expected0 := testutil.KC().WithCtxs(testutil.Ctx("c1")).ToYAML(t)
	if diff := cmp.Diff(expected0, out0); diff != "" {
		t.Fatalf("file 0 diff: %s", diff)
	}

	// Second file should have namespace set
	out1 := tl.OutputOf(1)
	expected1 := testutil.KC().WithCtxs(testutil.Ctx("c2").Ns("my-ns")).ToYAML(t)
	if diff := cmp.Diff(expected1, out1); diff != "" {
		t.Fatalf("file 1 diff: %s", diff)
	}
}
