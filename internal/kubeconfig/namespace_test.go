package kubeconfig

import (
	"testing"

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
