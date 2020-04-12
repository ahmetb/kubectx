package kubeconfig

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	err := new(Kubeconfig).WithLoader(&testLoader{in: strings.NewReader(`a:b`)}).Parse()
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
