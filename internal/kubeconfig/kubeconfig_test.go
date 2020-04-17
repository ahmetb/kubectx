package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestParse(t *testing.T) {
	err := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`a: [1, 2`)).Parse()
	if err == nil {
		t.Fatal("expected error from bad yaml")
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`[1, 2, 3]`)).Parse()
	if err == nil {
		t.Fatal("expected error from not-mapping root node")
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`current-context: foo`)).Parse()
	if err != nil {
		t.Fatal(err)
	}

	err = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCurrentCtx("foo").
		WithCtxs().ToYAML(t))).Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSave(t *testing.T) {
	in := "a: [1, 2, 3]\n"
	test := WithMockKubeconfigLoader(in)
	kc := new(Kubeconfig).WithLoader(test)
	defer kc.Close()
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.ModifyCurrentContext("hello"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}
	expected := "a: [1, 2, 3]\ncurrent-context: hello\n"
	if diff := cmp.Diff(expected, test.Output()); diff != "" {
		t.Fatal(diff)
	}
}
