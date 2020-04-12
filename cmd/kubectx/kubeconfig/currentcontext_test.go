package kubeconfig

import (
	"strings"
	"testing"
)

func TestKubeconfig_GetCurrentContext(t *testing.T) {
	tl := &testLoader{in: strings.NewReader(`current-context: foo`)}
	kc := new(Kubeconfig).WithLoader(tl)
	_, err := kc.ParseRaw()
	if err != nil {
		t.Fatal(err)
	}
	v := kc.GetCurrentContext()

	expected := "foo"
	if v != expected {
		t.Fatalf("expected=%q; got=%q", expected, v)
	}
}

func TestKubeconfig_GetCurrentContext_missingField(t *testing.T) {
	tl := &testLoader{in: strings.NewReader(`abc: def`)}
	kc := new(Kubeconfig).WithLoader(tl)
	_, err := kc.ParseRaw()
	if err != nil {
		t.Fatal(err)
	}
	v := kc.GetCurrentContext()

	expected := ""
	if v != expected {
		t.Fatalf("expected=%q; got=%q", expected, v)
	}
}

func TestKubeconfig_UnsetCurrentContext(t *testing.T) {
	tl := &testLoader{in: strings.NewReader(`current-context: foo`)}
	kc := new(Kubeconfig).WithLoader(tl)
	_, err := kc.ParseRaw()
	if err != nil {
		t.Fatal(err)
	}
	if err := kc.UnsetCurrentContext(); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	out := tl.out.String()
	expected := `current-context: ""
`
	if out != expected {
		t.Fatalf("expected=%q; got=%q", expected, out)
	}
}
