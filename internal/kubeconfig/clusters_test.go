package kubeconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func TestKubeconfig_ClusterOfContext_ctxNotFound(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(testutil.Ctx("c1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	_, err := kc.ClusterOfContext("c2")
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestKubeconfig_ClusterOfContext(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1").Cluster("c1c1"),
			testutil.Ctx("c2").Cluster("c2c2")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	v1, err := kc.ClusterOfContext("c1")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := `c1c1`; v1 != expected {
		t.Fatalf("c1: expected=\"%s\" got=\"%s\"", expected, v1)
	}

	v2, err := kc.ClusterOfContext("c2")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := `c2c2`; v2 != expected {
		t.Fatalf("c2: expected=\"%s\" got=\"%s\"", expected, v2)
	}
}

func TestKubeconfig_DeleteClusterEntry_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`[1, 2, 3]`))
	_ = kc.Parse()
	err := kc.DeleteClusterEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail on non-mapping nodes")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`a: b`))
	_ = kc.Parse()
	err = kc.DeleteClusterEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if clusters key does not exist")
	}

	kc = new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(`clusters: "some string"`))
	_ = kc.Parse()
	err = kc.DeleteClusterEntry("foo")
	if err == nil {
		t.Fatal("supposed to fail if clusters key is not an array")
	}
}

func TestKubeconfig_DeleteClusterEntry(t *testing.T) {
	test := WithMockKubeconfigLoader(
		testutil.KC().WithClusters(
			testutil.Cluster("c1"),
			testutil.Cluster("c2"),
			testutil.Cluster("c3")).ToYAML(t))
	kc := new(Kubeconfig).WithLoader(test)
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}
	if err := kc.DeleteClusterEntry("c1"); err != nil {
		t.Fatal(err)
	}
	if err := kc.Save(); err != nil {
		t.Fatal(err)
	}

	expected := testutil.KC().WithClusters(
		testutil.Cluster("c2"),
		testutil.Cluster("c3")).ToYAML(t)
	out := test.Output()
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("diff: %s", diff)
	}
}

func TestKubeconfig_CountClusterReferences_errors(t *testing.T) {
	kc := new(Kubeconfig).WithLoader(WithMockKubeconfigLoader(testutil.KC().
		WithCtxs(
			testutil.Ctx("c1").Cluster("c1c1"),
			testutil.Ctx("c2").Cluster("c2c2"),
			testutil.Ctx("c3").Cluster("c1c1")).ToYAML(t)))
	if err := kc.Parse(); err != nil {
		t.Fatal(err)
	}

	count1, err := kc.CountClusterReferences("c1c1")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := 2; count1 != expected {
		t.Fatalf("c1: expected=\"%d\" got=\"%d\"", expected, count1)
	}

	count2, err := kc.CountClusterReferences("c2c2")
	if err != nil {
		t.Fatal("unexpected err", err)
	}
	if expected := 1; count2 != expected {
		t.Fatalf("c1: expected=\"%d\" got=\"%d\"", expected, count2)
	}
}
