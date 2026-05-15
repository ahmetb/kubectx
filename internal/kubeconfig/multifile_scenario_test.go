package kubeconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSetNamespace_MultiFileScenario tests the scenario where current-context is in
// file1 but the actual context definition is in file2.
func TestSetNamespace_MultiFileScenario(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "config1")
	f2 := filepath.Join(dir, "config2")

	// File 1: has current-context but the context definition is in file 2
	os.WriteFile(f1, []byte(`apiVersion: v1
kind: Config
current-context: ctx-in-file2
contexts: []
clusters: []
users: []
`), 0644)

	// File 2: has the actual context definition
	os.WriteFile(f2, []byte(`apiVersion: v1
kind: Config
contexts:
- context:
    cluster: my-cluster
    user: my-user
  name: ctx-in-file2
clusters:
- cluster:
    server: https://test
  name: my-cluster
users:
- name: my-user
  user: {}
`), 0644)

	t.Setenv("KUBECONFIG", f1+":"+f2)

	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	cur, err := kc.GetCurrentContext()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Current context: %s", cur)

	if err := kc.SetNamespace("ctx-in-file2", "my-namespace"); err != nil {
		t.Fatalf("SetNamespace failed: %v", err)
	}

	if err := kc.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	c1, _ := os.ReadFile(f1)
	c2, _ := os.ReadFile(f2)
	t.Logf("File1 after save:\n%s", c1)
	t.Logf("File2 after save:\n%s", c2)

	if !strings.Contains(string(c2), "namespace: my-namespace") {
		t.Errorf("FAIL: namespace NOT saved to file2. File2:\n%s", c2)
	}
}

// TestSetNamespace_CurrentContextInFile2 tests when the current-context is defined
// in file2 (where kubectl logs "Config loaded from file: /path/to/config" for multiple files,
// it shows the first file, potentially misleading).
func TestSetNamespace_AllInSameFile(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "config1")

	os.WriteFile(f1, []byte(`apiVersion: v1
kind: Config
current-context: my-ctx
contexts:
- context:
    cluster: my-cluster
    user: my-user
  name: my-ctx
clusters:
- cluster:
    server: https://test
  name: my-cluster
users:
- name: my-user
  user: {}
`), 0644)

	t.Setenv("KUBECONFIG", f1)

	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if err := kc.SetNamespace("my-ctx", "my-namespace"); err != nil {
		t.Fatalf("SetNamespace failed: %v", err)
	}
	if err := kc.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	content, _ := os.ReadFile(f1)
	t.Logf("File after save:\n%s", content)

	if !strings.Contains(string(content), "namespace: my-namespace") {
		t.Errorf("FAIL: namespace NOT in file. Content:\n%s", content)
	}
}
