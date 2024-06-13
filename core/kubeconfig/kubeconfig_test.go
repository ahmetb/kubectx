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

	"github.com/ahmetb/kubectx/core/testutil"
)

func TestParse(t *testing.T) {
	err := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`a: [1, 2`)).Parse()
	if err == nil {
		t.Fatal("expected error from bad yaml")
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`[1, 2, 3]`)).Parse()
	if err == nil {
		t.Fatal("expected error from not-mapping root node")
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`current-context: foo`)).Parse()
	if err != nil {
		t.Fatal(err)
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCurrentCtx("foo").
		WithCtxs().ToYAML(t))).Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	in := "a: [1, 2, 3]\n"
	test := WithMockKubeconfigLoader(in)
	kc := new(Kubeconfig).WithLoader(test)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("hello"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}
	expected := "a: [1, 2, 3]\ncurrent-context: hello\n"
	if diff := cmp.Diff(expected, test.Output()); diff != "" {
		t.Fatal(diff)
	}
}
