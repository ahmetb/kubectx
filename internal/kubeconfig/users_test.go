package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_UserOfContext_ctxNotFound(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(testutil.Ctx("c1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	_, err := kc.UserOfContext("c2")
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestKubeconfig_UserOfContext(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1").User("c1u1"),
			testutil.Ctx("c2").User("c2u2")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	v1, err := kc.UserOfContext("c1")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := `c1u1`; v1 != expected {
		t.Fatalf("c1: expected=\"%s\" got=\"%s\"", expected, v1)
	}

	v2, err := kc.UserOfContext("c2")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := `c2u2`; v2 != expected {
		t.Fatalf("c2: expected=\"%s\" got=\"%s\"", expected, v2)
	}
}

func TestKubeconfig_DeleteUserEntry_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`[1, 2, 3]`))
	_ = kc.Parse()
	err := kc.DeleteUserEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail on non-mapping nodes")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`a: b`))
	_ = kc.Parse()
	err = kc.DeleteUserEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if users key does not exist")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`users: "some string"`))
	_ = kc.Parse()
	err = kc.DeleteUserEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if users key is not an array")
	}
}

func TestKubeconfig_DeleteUserEntry(t *testing.T) {
	test := WithMockKubeconfigLoader(
		testutil.KC().WithUsers(
			testutil.User("u1"),
			testutil.User("u2"),
			testutil.User("u3")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.DeleteUserEntry("u1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithUsers(
		testutil.User("u2"),
		testutil.User("u3")).ToYAML(t)
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_CountUserReferences_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1").User("c1u1"),
			testutil.Ctx("c2").User("c2u2"),
			testutil.Ctx("c3").User("c1u1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	count1, err := kc.CountUserReferences("c1u1")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := 2; count1 != expected {
		t.Fatalf("c1: expected=\"%d\" got=\"%d\"", expected, count1)
	}

	count2, err := kc.CountUserReferences("c2u2")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := 1; count2 != expected {
		t.Fatalf("c1: expected=\"%d\" got=\"%d\"", expected, count2)
	}
}
