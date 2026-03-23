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
	v, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatal(err)
	}

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
	v, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatal(err)
	}

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

func TestKubeconfig_GetCurrentContext_MultiFile_FirstNonEmpty(t *testing.T) {
	cfg1 := testutil.KC().WithCtxs(testutil.Ctx("ctx1")).ToYAML(t) // no current-context
	cfg2 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	v, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatal(err)
	}
	if v != "ctx2" {
		t.Fatalf("expected=\"ctx2\"; got=\"%s\"", v)
	}
}

func TestKubeconfig_GetCurrentContext_MultiFile_FirstWins(t *testing.T) {
	cfg1 := testutil.KC().WithCurrentCtx("ctx1").WithCtxs(testutil.Ctx("ctx1")).ToYAML(t)
	cfg2 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	v, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatal(err)
	}
	if v != "ctx1" {
		t.Fatalf("expected=\"ctx1\"; got=\"%s\"", v)
	}
}

func TestKubeconfig_UnsetCurrentContext_MultiFile(t *testing.T) {
	cfg1 := testutil.KC().WithCurrentCtx("ctx1").WithCtxs(testutil.Ctx("ctx1")).ToYAML(t)
	cfg2 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	tl := WithMockMultiKubeconfigLoader(cfg1, cfg2)
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

	// After unsetting, GetCurrentContext should return ctx2 (from second file)
	// Re-parse the saved output to verify
	out0 := tl.OutputOf(0)
	expected0 := testutil.KC().WithCurrentCtx("").WithCtxs(testutil.Ctx("ctx1")).ToYAML(t)
	if out0 != expected0 {
		t.Fatalf("file 0: expected=\"%s\"; got=\"%s\"", expected0, out0)
	}

	// Second file should be unchanged
	out1 := tl.OutputOf(1)
	expected1 := testutil.KC().WithCurrentCtx("ctx2").WithCtxs(testutil.Ctx("ctx2")).ToYAML(t)
	if out1 != expected1 {
		t.Fatalf("file 1: expected=\"%s\"; got=\"%s\"", expected1, out1)
	}
}
