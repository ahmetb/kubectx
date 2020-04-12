package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_kubeconfigPath_homePath(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", "/foo/bar")
	defer os.Setenv("HOME", origHome)

	got, err := kubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "config")

	if got != expected{
		t.Fatalf("wrong value: expected=%s got=%s", expected, got)
	}
}

func Test_kubeconfigPath_userprofile(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Unsetenv("HOME")
	os.Setenv("USERPROFILE", "/foo/bar")
	defer os.Setenv("HOME", origHome)

	got, err := kubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "config")

	if got != expected{
		t.Fatalf("wrong value: expected=%s got=%s", expected, got)
	}
}

func Test_kubeconfigPath_noEnvVars(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserprofile := os.Getenv("USERPROFILE")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer os.Setenv("HOME", origHome)
	defer os.Setenv("USERPROFILE", origUserprofile)

	_, err := kubeconfigPath()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_kubeconfigPath_envOvveride(t *testing.T) {
	os.Setenv("KUBECONFIG", "foo")
	defer os.Unsetenv("KUBECONFIG")

	v, err := kubeconfigPath()
	if err != nil { t.Fatal(err)}
	if expected := "foo"; v != expected {
		t.Fatalf("expected=%q, got=%q", expected, v)
	}
}

func Test_kubeconfigPath_envOvverideDoesNotSupportPathSeparator(t *testing.T) {
	path := strings.Join([]string{"file1","file2"}, string(os.PathListSeparator))
	os.Setenv("KUBECONFIG", path)
	defer os.Unsetenv("KUBECONFIG")

	_, err := kubeconfigPath()
	if err == nil { t.Fatal("expected error")}
}

func testfile(t *testing.T, contents string) (path string, cleanup func()) {
	t.Helper()

	f, err := ioutil.TempFile(os.TempDir(), "test-file")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	path = f.Name()
	if _, err := f.Write([]byte(contents)); err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}

	return path, func() {
		f.Close()
		os.Remove(path)
	}
}
