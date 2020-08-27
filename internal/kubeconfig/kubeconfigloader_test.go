package kubeconfig

import (
	"github.com/ahmetb/kubectx/internal/cmdutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func Test_kubeconfigPath(t *testing.T) {
	defer testutil.WithEnvVar("HOME", "/x/y/z")()

	expected := filepath.FromSlash("/x/y/z/.kube/config")
	got, err := FindKubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Fatalf("got=%q expected=%q", got, expected)
	}
}

func Test_kubeconfigPath_noEnvVars(t *testing.T) {
	defer testutil.WithEnvVar("XDG_CACHE_HOME", "")()
	defer testutil.WithEnvVar("HOME", "")()
	defer testutil.WithEnvVar("USERPROFILE", "")()

	_, err := FindKubeconfigPath()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_kubeconfigPath_envOverride(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()

	v, err := FindKubeconfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if expected := "foo"; v != expected {
		t.Fatalf("expected=%q, got=%q", expected, v)
	}
}

func Test_kubeconfigPath_envOverrideDoesNotSupportPathSeparator(t *testing.T) {
	path := strings.Join([]string{"file1", "file2"}, string(os.PathListSeparator))
	defer testutil.WithEnvVar("KUBECONFIG", path)()

	_, err := FindKubeconfigPath()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStandardKubeconfigLoader_returnsNotFoundErr(t *testing.T) {
	defer testutil.WithEnvVar("KUBECONFIG", "foo")()
	kc := new(Kubeconfig).WithLoader(DefaultLoader)
	err := kc.Parse()
	if err == nil {
		t.Fatal("expected err")
	}
	if !cmdutil.IsNotFoundErr(err) {
		t.Fatalf("expected ENOENT error; got=%v", err)
	}
}
