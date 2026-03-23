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

func TestKubeconfig_DeleteContextEntry_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`[1, 2, 3]`))
	_ = kc.Parse()
	err := kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail on non-mapping nodes")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`a: b`))
	_ = kc.Parse()
	err = kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if contexts key does not exist")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`contexts: "some string"`))
	_ = kc.Parse()
	err = kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if contexts key is not an array")
	}
}

func TestKubeconfig_DeleteContextEntry(t *testing.T) {
	test := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2"),
			testutil.Ctx("c3")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.DeleteContextEntry("c1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithCtxs(
		testutil.Ctx("c2"),
		testutil.Ctx("c3")).ToYAML(t)
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyCurrentContext_fieldExists(t *testing.T) {
	test := WithMockKubeconfigLoader(
		testutil.KC().WithCurrentCtx("abc").Set("field1", "value1").ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("foo"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithCurrentCtx("foo").Set("field1", "value1").ToYAML(t)
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyCurrentContext_fieldMissing(t *testing.T) {
	test := WithMockKubeconfigLoader(`f1: v1`)
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("foo"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := `f1: v1
current-context: foo
`
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyContextName_noContextsEntryError(t *testing.T) {
	// no context entries
	test := WithMockKubeconfigLoader(`a: b`)
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c1", "c2"); err == nil {
		t.Fatal("was expecting error for no 'contexts' entry; got nil")
	}
}

func TestKubeconfig_ModifyContextName_contextsEntryNotSequenceError(t *testing.T) {
	// no context entries
	test := WithMockKubeconfigLoader(
		`contexts: "hello"`)
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c1", "c2"); err == nil {
		t.Fatal("was expecting error for 'context entry not a sequence'; got nil")
	}
}

func TestKubeconfig_ModifyContextName_noChange(t *testing.T) {
	test := WithMockKubeconfigLoader(testutil.KC().WithCtxs(
		testutil.Ctx("c1"),
		testutil.Ctx("c2"),
		testutil.Ctx("c3")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c5", "c6"); err == nil {
		t.Fatal("was expecting error for 'no changes made'")
	}
}

func TestKubeconfig_ModifyContextName(t *testing.T) {
	test := WithMockKubeconfigLoader(testutil.KC().WithCtxs(
		testutil.Ctx("c1"),
		testutil.Ctx("c2"),
		testutil.Ctx("c3")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c1", "ccc"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithCtxs(
		testutil.Ctx("ccc"),
		testutil.Ctx("c2"),
		testutil.Ctx("c3")).ToYAML(t)
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyCurrentContext_MultiFile_WritesToFirst(t *testing.T) {
	cfg1 := testutil.KC().WithCurrentCtx("ctx1").WithCtxs(testutil.Ctx("ctx1")).ToYAML(t)
	cfg2 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("ctx2"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	// First file should have new current-context
	out0 := tl.OutputOf(0)
	expected0 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx1")).ToYAML(t)
	if diff := cmp.Diff(expected0, out0); diff != "" {
		t.Fatalf("file 0 diff: %s", diff)
	}

	// Second file should be unchanged
	out1 := tl.OutputOf(1)
	expected1 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	if diff := cmp.Diff(expected1, out1); diff != "" {
		t.Fatalf("file 1 diff: %s", diff)
	}
}

func TestKubeconfig_DeleteContextEntry_MultiFile_FromCorrectFile(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("c1"), testutil.Ctx("c2")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c3"), testutil.Ctx("c4")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	// Delete c3 which is in file 2
	if err := kc.DeleteContextEntry("c3"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	// First file should be unchanged
	out0 := tl.OutputOf(0)
	expected0 := testutil.KC().WithCtxs(testutil.Ctx("c1"), testutil.Ctx("c2")).ToYAML(t)
	if diff := cmp.Diff(expected0, out0); diff != "" {
		t.Fatalf("file 0 diff: %s", diff)
	}

	// Second file should have c3 removed
	out1 := tl.OutputOf(1)
	expected1 := testutil.KC().WithCtxs(testutil.Ctx("c4")).ToYAML(t)
	if diff := cmp.Diff(expected1, out1); diff != "" {
		t.Fatalf("file 1 diff: %s", diff)
	}
}

func TestKubeconfig_ModifyContextName_MultiFile_InCorrectFile(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("c1")).ToYAML(t)
	cfg2 := testutil.KC().WithCtxs(testutil.Ctx("c2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	// Rename c2 (in file 2) to c2-new
	if err := kc.ModifyContextName("c2", "c2-new"); err != nil {
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

	// Second file should have c2 renamed to c2-new
	out1 := tl.OutputOf(1)
	expected1 := testutil.KC().WithCtxs(testutil.Ctx("c2-new")).ToYAML(t)
	if diff := cmp.Diff(expected1, out1); diff != "" {
		t.Fatalf("file 1 diff: %s", diff)
	}
}
