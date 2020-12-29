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
