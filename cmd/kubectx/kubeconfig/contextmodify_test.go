package kubeconfig

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestKubeconfig_DeleteContextEntry_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`[1, 2, 3]`)})
	_ = kc.Parse()
	err := kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail on non-mapping nodes")
	}

	kc = new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`a: b`)})
	_ = kc.Parse()
	err = kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if contexts key does not exist")
	}

	kc = new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`contexts: "some string"`)})
	_ = kc.Parse()
	err = kc.DeleteContextEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if contexts key is not an array")
	}
}

func TestKubeconfig_DeleteContextEntry(t *testing.T) {
	test := &testLoader{in: strings.NewReader(
		`contexts: [{name: c1}, {name: c2}, {name: c3}]`)}
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.DeleteContextEntry("c1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := "contexts: [{name: c2}, {name: c3}]\n"
	out := test.out.String()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyCurrentContext_fieldExists(t *testing.T) {
	test := &testLoader{in: strings.NewReader(
		`current-context: abc
field1: value1`)}
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("foo"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := `current-context: foo
field1: value1` + "\n"
	out := test.out.String()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyCurrentContext_fieldMissing(t *testing.T) {
	test := &testLoader{in: strings.NewReader(
		`field1: value1`)}
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("foo"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := `field1: value1` + "\n" + "current-context: foo\n"
	out := test.out.String()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_ModifyContextName_noChange(t *testing.T) {
	test := &testLoader{in: strings.NewReader(
		`contexts: [{name: c1}, {name: c2}, {name: c3}]`)}
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c5", "c6"); err == nil {
		t.Fatal("was expecting error for 'no changes made'")
	}
}

func TestKubeconfig_ModifyContextName(t *testing.T) {
	test := &testLoader{in: strings.NewReader(
		`contexts: [{name: c1}, {name: c2}, {name: c3}]`)}
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyContextName("c1", "ccc"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := "contexts: [{name: ccc}, {name: c2}, {name: c3}]\n"
	out := test.out.String()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}
