package main

import (
	"io/ioutil"
	"os"
	"testing"
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
		t.Fatalf("Load()=%q; expected=%q", v, expected)
	}
}
