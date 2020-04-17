package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_ContextNames(t *testing.T) {
	tl := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("abc"),
			testutil.Ctx("def"),
			testutil.Ctx("ghi")).Set("field1", map[string]string{"bar": "zoo"}).ToYAML(t))
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

func TestKubeconfig_ContextNames_noContextsEntry(t *testing.T) {
	tl := WithMockKubeconfigLoader(`a: b`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	ctx := kc.ContextNames()
	var expected []string = nil
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_ContextNames_nonArrayContextsEntry(t *testing.T) {
	tl := WithMockKubeconfigLoader(`contexts: "hello"`)
	kc := new(Kubeconfig).WithLoader(tl)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	ctx := kc.ContextNames()
	var expected []string = nil
	if diff := cmp.Diff(expected, ctx); diff != "" {
		t.Fatalf("%s", diff)
	}
}

func TestKubeconfig_CheckContextExists(t *testing.T) {
	tl := WithMockKubeconfigLoader(
		testutil.KC().WithCtxs(
			testutil.Ctx("c1"),
			testutil.Ctx("c2")).ToYAML(t))

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
