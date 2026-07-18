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

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

func writeTestKubeconfigWithNamespace(t *testing.T) string {
	t.Helper()
	const cfg = `apiVersion: v1
kind: Config
current-context: ctx-a
contexts:
- name: ctx-a
  context: { cluster: cluster-a, user: user-a, namespace: ns1 }
clusters:
- name: cluster-a
  cluster: { server: https://example.invalid }
users:
- name: user-a
  user: { token: fake }
`
	p := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(p, []byte(cfg), 0644); err != nil {
		t.Fatalf("write kubeconfig: %v", err)
	}
	return p
}

func installKubensFakeFzf(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fake fzf shell script unsupported on windows")
	}
	binDir := t.TempDir()
	const script = `#!/bin/sh
cat > "$KUBECTX_TEST_FZF_STDIN"
printf '%s\n' "$KUBECTX_TEST_FZF_OUT"
`
	if err := os.WriteFile(filepath.Join(binDir, "fzf"), []byte(script), 0755); err != nil {
		t.Fatalf("write fake fzf: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func Test_formatNamespaceList(t *testing.T) {
	t.Setenv("KUBECONFIG", writeTestKubeconfigWithNamespace(t))
	t.Setenv("_MOCK_NAMESPACES", "1")

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatalf("parse: %v", err)
	}

	var out bytes.Buffer
	if err := formatNamespaceList(kc, &out); err != nil {
		t.Fatalf("formatNamespaceList: %v", err)
	}
	got := out.String()
	for _, want := range []string{"ns1", "ns2"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected list to contain %q; got %q", want, got)
		}
	}
}

func TestInteractiveSwitchOp_pipesNamespaceListToFzfStdin(t *testing.T) {
	t.Setenv("KUBECONFIG", writeTestKubeconfigWithNamespace(t))
	t.Setenv("_MOCK_NAMESPACES", "1")
	installKubensFakeFzf(t)

	stdinRec := filepath.Join(t.TempDir(), "fzf-stdin.txt")
	t.Setenv("KUBECTX_TEST_FZF_STDIN", stdinRec)
	// pick ns2 (not the current ns1) so we exercise the switch path
	t.Setenv("KUBECTX_TEST_FZF_OUT", "ns2")

	var stderr bytes.Buffer
	if err := (InteractiveSwitchOp{}).Run(io.Discard, &stderr); err != nil {
		t.Fatalf("InteractiveSwitchOp.Run: %v", err)
	}

	rec, err := os.ReadFile(stdinRec)
	if err != nil {
		t.Fatalf("read fzf stdin capture: %v", err)
	}
	if !strings.Contains(string(rec), "ns1") || !strings.Contains(string(rec), "ns2") {
		t.Fatalf("fzf stdin did not contain namespace list; got %q", string(rec))
	}

	// switchNamespace should have set ns2 as the active namespace
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	ns, err := kc.NamespaceOfContext("ctx-a")
	if err != nil {
		t.Fatalf("namespace of context: %v", err)
	}
	if ns != "ns2" {
		t.Errorf("active namespace = %q, want %q", ns, "ns2")
	}
}
