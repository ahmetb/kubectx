package kubeconfig

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	err := new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`a: [1, 2`)}).Parse()
	if err == nil {
		t.Fatal("expected error from bad yaml")
	}

	err = new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`[1, 2, 3]`)}).Parse()
	if err == nil {
		t.Fatal("expected error from not-mapping root node")
	}

	err = new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`current-context: foo`)}).Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	in := `a: [1, 2, 3]` + "\n"
	test := &testLoader{in: strings.NewReader(in)}
	kc := new(Kubeconfig).WithLoader(test)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("hello"); err != nil {t.Fatal(err)}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}
	expected := in+"current-context: hello\n"
	if diff := cmp.Diff(expected, test.out.String()); diff != "" {
		t.Fatal(diff)
	}
}
