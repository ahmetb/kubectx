package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ahmetb/kubectx/internal/testutil"
)

var (
	tr, fa = true, false
)

func Test_useColors_forceColors(t *testing.T) {
	defer testutil.WithEnvVar("_KUBECTX_FORCE_COLOR", "1")()
	defer testutil.WithEnvVar("NO_COLOR", "1")()

	if v := UseColors(); !cmp.Equal(v, &tr) {
		t.Fatalf("expected UseColors() = true; got = %v", v)
	}
}

func Test_useColors_disableColors(t *testing.T) {
	defer testutil.WithEnvVar("NO_COLOR", "1")()

	if v := UseColors(); !cmp.Equal(v, &fa) {
		t.Fatalf("expected UseColors() = false; got = %v", v)
	}
}

func Test_useColors_default(t *testing.T) {
	defer testutil.WithEnvVar("NO_COLOR", "")()
	defer testutil.WithEnvVar("_KUBECTX_FORCE_COLOR", "")()

	if v := UseColors(); v != nil {
		t.Fatalf("expected UseColors() = nil; got=%v", *v)
	}
}
