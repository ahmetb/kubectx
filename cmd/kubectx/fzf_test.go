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

// writeTestKubeconfig writes a kubeconfig with contexts ctx-a (current)
// and ctx-b into a temp dir and returns its path.
func writeTestKubeconfig(t *testing.T) string {
	t.Helper()
	const cfg = `apiVersion: v1
kind: Config
current-context: ctx-a
contexts:
- name: ctx-a
  context: { cluster: cluster-a, user: user-a }
- name: ctx-b
  context: { cluster: cluster-b, user: user-b }
clusters:
- name: cluster-a
  cluster: { server: https://example.invalid }
- name: cluster-b
  cluster: { server: https://example.invalid }
users:
- name: user-a
  user: { token: fake }
- name: user-b
  user: { token: fake }
`
	p := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(p, []byte(cfg), 0644); err != nil {
		t.Fatalf("write kubeconfig: %v", err)
	}
	return p
}

// installFakeFzf installs a fake `fzf` in a temp bin dir that copies its
// stdin to the file at $KUBECTX_TEST_FZF_STDIN and echoes $KUBECTX_TEST_FZF_OUT
// to stdout. It prepends the dir to PATH. The test reads the recorded file to
// verify the candidate list was piped to fzf stdin.
func installFakeFzf(t *testing.T) {
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

func Test_formatContextList(t *testing.T) {
	t.Setenv("KUBECONFIG", writeTestKubeconfig(t))

	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatalf("parse: %v", err)
	}

	var out bytes.Buffer
	if err := formatContextList(kc, &out); err != nil {
		t.Fatalf("formatContextList: %v", err)
	}
	got := out.String()
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), got)
	}
	for _, l := range lines {
		// strip any ANSI color escapes from the active-context line
		name := strings.Trim(l, "\x1b[0;m123456789")
		if name != "ctx-a" && name != "ctx-b" {
			t.Errorf("unexpected line %q", l)
		}
	}
}

func TestInteractiveSwitchOp_pipesListToFzfStdin(t *testing.T) {
	cfg := writeTestKubeconfig(t)
	t.Setenv("KUBECONFIG", cfg)
	installFakeFzf(t)

	stdinRec := filepath.Join(t.TempDir(), "fzf-stdin.txt")
	t.Setenv("KUBECTX_TEST_FZF_STDIN", stdinRec)
	t.Setenv("KUBECTX_TEST_FZF_OUT", "ctx-b")

	var stderr bytes.Buffer
	if err := (InteractiveSwitchOp{}).Run(io.Discard, &stderr); err != nil {
		t.Fatalf("InteractiveSwitchOp.Run: %v", err)
	}

	rec, err := os.ReadFile(stdinRec)
	if err != nil {
		t.Fatalf("read fzf stdin capture: %v", err)
	}
	if !strings.Contains(string(rec), "ctx-a") || !strings.Contains(string(rec), "ctx-b") {
		t.Fatalf("fzf stdin did not contain context list; got %q", string(rec))
	}

	// switchContext should have switched to ctx-b
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	cur, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatalf("get current: %v", err)
	}
	if cur != "ctx-b" {
		t.Errorf("current-context = %q, want %q", cur, "ctx-b")
	}
}

func TestInteractiveDeleteOp_pipesListToFzfStdin(t *testing.T) {
	cfg := writeTestKubeconfig(t)
	t.Setenv("KUBECONFIG", cfg)
	installFakeFzf(t)

	stdinRec := filepath.Join(t.TempDir(), "fzf-stdin.txt")
	t.Setenv("KUBECTX_TEST_FZF_STDIN", stdinRec)
	t.Setenv("KUBECTX_TEST_FZF_OUT", "ctx-b")

	var stderr bytes.Buffer
	if err := (InteractiveDeleteOp{}).Run(io.Discard, &stderr); err != nil {
		t.Fatalf("InteractiveDeleteOp.Run: %v", err)
	}

	rec, err := os.ReadFile(stdinRec)
	if err != nil {
		t.Fatalf("read fzf stdin capture: %v", err)
	}
	if !strings.Contains(string(rec), "ctx-a") || !strings.Contains(string(rec), "ctx-b") {
		t.Fatalf("fzf stdin did not contain context list; got %q", string(rec))
	}

	// deleteContext should have removed ctx-b
	kc := new(kubeconfig.Kubeconfig).WithLoader(kubeconfig.DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	names, err := kc.ContextNames()
	if err != nil {
		t.Fatalf("context names: %v", err)
	}
	for _, n := range names {
		if n == "ctx-b" {
			t.Errorf("ctx-b should have been deleted; remaining contexts: %v", names)
		}
	}
	if len(names) != 1 || names[0] != "ctx-a" {
		t.Errorf("expected only ctx-a to remain; got %v", names)
	}
}
