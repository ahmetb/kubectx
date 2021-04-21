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
	"github.com/ahmetb/kubectx/internal/cmdutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func Test_kubeconfigPath(t *testing.T) {
	defer testutil.WithEnvVar("HOME", "/x/y/z")()

	expected := filepath.FromSlash("/x/y/z/.kube/config")
	got, err := kubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Fatalf("got=%q expected=%q", got, expected)
	}
}

func Test_kubeconfigPath_noEnvVars(t *testing.T) {
	defer testutil.WithEnvVar("XDG_CACHE_HOME", "")()
	defer testutil.WithEnvVar("HOME", "")()
	defer testutil.WithEnvVar("USERPROFILE", "")()

	_, err := kubeconfigPath()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_kubeconfigPath_envOvveride(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()

	v, err := kubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "foo"; v != expected {
		t.Fatalf("expected=%q, got=%q", expected, v)
	}
}

func Test_kubeconfigPath_envOvverideDoesNotSupportPathSeparator(t *testing.T) {
	path := strings.Join([]string{"file1", "file2"}, string(os.PathListSeparator))
	defer testutil.WithEnvVar("KUBECONFIG", path)()

	_, err := kubeconfigPath()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStandardKubeconfigLoader_returnsNotFoundErr(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()
	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	err := kc.Parse()
	if err == nil {
		t.Fatal("expected err")
	}
	if !cmdutil.IsNotFoundErr(err) {
		t.Fatalf("expected ENOENT error; got=%v", err)
	}
}
