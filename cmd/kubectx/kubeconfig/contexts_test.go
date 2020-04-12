package kubeconfig

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestKubeconfig_ContextNames(t *testing.T) {
	tl := &testLoader{in: strings.NewReader(`
contexts:
- name: abc
- name: def
  field1: value1
- name: ghi
  foo:
    bar: zoo`)}

	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	ctx := kc.ContextNames()
	expected := []string{"abc", "def", "ghi"}
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_CheckContextExists(t *testing.T) {
	tl := &testLoader{in: strings.NewReader(`contexts:
- name: c1
- name: c2`)}

	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if !kc.ContextExists("c1") {
		t.Fatal("c1 actually exists; reported false")
	}
	if !kc.ContextExists("c2") {
		t.Fatal("c2 actually exists; reported false")
	}
	if kc.ContextExists("c3") {
		t.Fatal("c3 does not exist; but reported true")
	}
}
