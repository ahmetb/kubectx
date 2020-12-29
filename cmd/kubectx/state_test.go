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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
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
	path, cleanup := testutil.TempFile(t, "foo")
	defer cleanup()

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
		t.Fatalf("read wrong value=\"%s\"; expected=\"%s\"", v, expected)
	}
}

func Test_kubectxFilePath(t *testing.T) {
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.FromSlash("/foo/bar"))
	defer os.Setenv("HOME", origHome)

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
	origHome := os.Getenv("HOME")
	origUserprofile := os.Getenv("USERPROFILE")
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")
	defer os.Setenv("HOME", origHome)
	defer os.Setenv("USERPROFILE", origUserprofile)

	_, err := kubectxPrevCtxFile()
	if err == nil {
		t.Fatal(err)
	}
}
