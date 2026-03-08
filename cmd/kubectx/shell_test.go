package main

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/ahmetb/kubectx/internal/env"
)

func Test_detectShell_unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping unix shell detection test on windows")
	}

	tests := []struct {
		name     string
		shellEnv string
		want     string
	}{
		{
			name:     "SHELL env set",
			shellEnv: "/bin/zsh",
			want:     "/bin/zsh",
		},
		{
			name:     "SHELL env empty, falls back to /bin/sh",
			shellEnv: "",
			want:     "/bin/sh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := os.Getenv("SHELL")
			defer os.Setenv("SHELL", orig)

			os.Setenv("SHELL", tt.shellEnv)
			if tt.shellEnv == "" {
				os.Unsetenv("SHELL")
			}

			got := detectShell()
			if got != tt.want {
				t.Errorf("detectShell() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_ShellOp_blockedWhenNested(t *testing.T) {
	// Simulate being inside an isolated shell
	orig := os.Getenv(env.EnvIsolatedShell)
	defer os.Setenv(env.EnvIsolatedShell, orig)
	os.Setenv(env.EnvIsolatedShell, "1")

	op := ShellOp{Target: "some-context"}
	var stdout, stderr bytes.Buffer
	err := op.Run(&stdout, &stderr)

	if err == nil {
		t.Fatal("expected error when running ShellOp inside isolated shell, got nil")
	}

	want := "locked single-context shell to"
	if !bytes.Contains([]byte(err.Error()), []byte(want)) {
		// The error may not contain the context name if kubeconfig is not available,
		// but it should still be blocked
		want2 := "locked single-context shell"
		if !bytes.Contains([]byte(err.Error()), []byte(want2)) {
			t.Errorf("error message %q does not contain %q", err.Error(), want2)
		}
	}
}

func Test_resolveKubectl_envVar(t *testing.T) {
	orig := os.Getenv("KUBECTL")
	defer os.Setenv("KUBECTL", orig)

	os.Setenv("KUBECTL", "/custom/path/kubectl")
	got, err := resolveKubectl()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom/path/kubectl" {
		t.Errorf("resolveKubectl() = %q, want %q", got, "/custom/path/kubectl")
	}
}

func Test_resolveKubectl_inPath(t *testing.T) {
	orig := os.Getenv("KUBECTL")
	defer os.Setenv("KUBECTL", orig)
	os.Unsetenv("KUBECTL")

	// kubectl should be findable in PATH on most dev machines
	got, err := resolveKubectl()
	if err != nil {
		t.Skip("kubectl not in PATH, skipping")
	}
	if got == "" {
		t.Error("resolveKubectl() returned empty string")
	}
}

func Test_checkIsolatedMode_notSet(t *testing.T) {
	orig := os.Getenv(env.EnvIsolatedShell)
	defer os.Setenv(env.EnvIsolatedShell, orig)
	os.Unsetenv(env.EnvIsolatedShell)

	err := checkIsolatedMode()
	if err != nil {
		t.Errorf("expected nil error when not in isolated mode, got: %v", err)
	}
}

func Test_checkIsolatedMode_set(t *testing.T) {
	orig := os.Getenv(env.EnvIsolatedShell)
	defer os.Setenv(env.EnvIsolatedShell, orig)
	os.Setenv(env.EnvIsolatedShell, "1")

	err := checkIsolatedMode()
	if err == nil {
		t.Fatal("expected error when in isolated mode, got nil")
	}

	want := "locked single-context shell"
	if !bytes.Contains([]byte(err.Error()), []byte(want)) {
		t.Errorf("error message %q does not contain %q", err.Error(), want)
	}
}
