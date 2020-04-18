package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TempFile(t *testing.T, contents string) (path string, cleanup func()) {
	// TODO consider removing, used only in one place.
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
