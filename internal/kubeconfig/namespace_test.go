package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_NamespaceOfContext_ctxNotFound(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(testutil.Ctx("c1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	_, err := kc.NamespaceOfContext("c2")
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestKubeconfig_NamespaceOfContext(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2").Ns("c2n1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	v1, err := kc.NamespaceOfContext("c1")
	if err != nil {
		t.Fatal("expected err")
	}
	if expected := `default`; v1 != expected {
		t.Fatalf("c1: expected=%q got=%q", expected, v1)
	}

	v2, err := kc.NamespaceOfContext("c2")
	if err != nil {
		t.Fatal("expected err")
	}
	if expected := `c2n1`; v2 != expected {
		t.Fatalf("c2: expected=%q got=%q", expected, v2)
	}
}

func TestKubeconfig_SetNamespace(t *testing.T) {
	l := WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2").Ns("c2n1")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(l)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	if err := kc.SetNamespace("c3", "foo"); err == nil {
		t.Fatalf("expected error for non-existing ctx")
	}

	if err := kc.SetNamespace("c1", "c1n1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.SetNamespace("c2", "c2n2"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithCtxs(
		testutil.Ctx("c1").Ns("c1n1"),
		testutil.Ctx("c2").Ns("c2n2")).ToYAML(t)
	if diff := cmp.Diff(l.Output(), expected); diff != "" {
		t.Fatal(diff)
	}
}
