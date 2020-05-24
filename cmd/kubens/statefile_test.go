package main

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestNSFile(t *testing.T) {
	td, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	f := NewNSFile("foo")
	f.dir = td
	v, err := f.Load()
	if err != nil {
		t.Fatal(err)
	}
	if v != "" {
		t.Fatalf("Load() expected empty; got=%v", err)
	}

	err = f.Save("bar")
	if err != nil {
		t.Fatalf("Save() err=%v", err)
	}

	v, err = f.Load()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "bar"; v != expected {
		t.Fatalf("Load()=\"%s\"; expected=\"%s\"", v, expected)
	}
}

func TestNSFile_path_windows(t *testing.T) {
	defer testutil.WithEnvVar("_FORCE_GOOS", "windows")()
	fp := NewNSFile("a:b:c").path()

	if expected := "a__b__c"; !strings.HasSuffix(fp, expected) {
		t.Fatalf("file did not have expected ending %q: %s", expected, fp)
	}
}

func Test_isWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("won't test this case on windows")
	}

	got := isWindows()
	if got {
		t.Fatalf("isWindows() returned true for %s", runtime.GOOS)
	}

	defer testutil.WithEnvVar("_FORCE_GOOS", "windows")()
	if !isWindows() {
		t.Fatalf("isWindows() failed to detect windows with env override.")
	}
}
