package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_readLastContext_nonExistingFile(t *testing.T) {
	s, err := readLastContext(filepath.FromSlash("/non/existing/file"))
	if err != nil {
		t.Fatal(err)
	}
	if s != "" {
		t.Fatalf("expected empty string; got=%q", s)
	}
}

func Test_readLastContext(t *testing.T) {
	path, cleanup := testfile(t, "foo")
	defer cleanup()

	s, err := readLastContext(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "foo"; s != expected {
		t.Fatalf("expected=%q; got=%q", expected, s)
	}
}

func Test_writeLastContext_err(t *testing.T) {
	path := filepath.Join(os.DevNull, "foo", "bar")
	err := writeLastContext(path, "foo")
	if err == nil {
		t.Fatal("got empty error")
	}
}

func Test_writeLastContext(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "state-file-test")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "foo", "bar")

	if err := writeLastContext(path, "ctx1"); err != nil {
		t.Fatal(err)
	}

	v, err := readLastContext(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "ctx1"; v != expected {
		t.Fatalf("read wrong value=%q; expected=%q", v, expected)
	}
}

func Test_kubectxFilePath(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.FromSlash("/foo/bar"))
	defer os.Setenv("HOME", origHome)

	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "kubectx")
	v, err := kubectxFilePath()
	if err != nil {
		t.Fatal(err)
	}
	if v != expected {
		t.Fatalf("expected=%q got=%q", expected, v)
	}
}

func Test_kubectxFilePath_error(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserprofile := os.Getenv("USERPROFILE")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer os.Setenv("HOME", origHome)
	defer os.Setenv("USERPROFILE", origUserprofile)

	_, err := kubectxFilePath()
	if err == nil {
		t.Fatal(err)
	}
}
