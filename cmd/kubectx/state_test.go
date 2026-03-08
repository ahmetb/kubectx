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
		t.Fatalf("expected empty string; got=\"%s\"", s)
	}
}

func Test_readLastContext(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")
	if err := os.WriteFile(path, []byte("foo"), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := readLastContext(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "foo"; s != expected {
		t.Fatalf("expected=\"%s\"; got=\"%s\"", expected, s)
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
	dir := t.TempDir()
	path := filepath.Join(dir, "foo", "bar")

	if err := writeLastContext(path, "ctx1"); err != nil {
		t.Fatal(err)
	}

	v, err := readLastContext(path)
	if err != nil {
		t.Fatal(err)
	}
	if expected := "ctx1"; v != expected {
		t.Fatalf("read wrong value=\"%s\"; expected=\"%s\"", v, expected)
	}
}

func Test_kubectxFilePath(t *testing.T) {
	t.Setenv("HOME", filepath.FromSlash("/foo/bar"))

	expected := filepath.Join(filepath.FromSlash("/foo/bar"), ".kube", "kubectx")
	v, err := kubectxPrevCtxFile()
	if err != nil {
		t.Fatal(err)
	}
	if v != expected {
		t.Fatalf("expected=\"%s\" got=\"%s\"", expected, v)
	}
}

func Test_kubectxFilePath_error(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	_, err := kubectxPrevCtxFile()
	if err == nil {
		t.Fatal(err)
	}
}
