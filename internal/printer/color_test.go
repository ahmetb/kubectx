package printer

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func withTestVar(key, value string) func() {
	orig, ok := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if ok {
			os.Setenv(key, orig)
		} else {
			os.Unsetenv(key)
		}
	}
}

var (
	tr, fa = true, false
)

func Test_useColors_forceColors(t *testing.T) {
	defer withTestVar("_KUBECTX_FORCE_COLOR", "1")()
	defer withTestVar("NO_COLOR", "1")()

	if v := UseColors(); !cmp.Equal(v, &tr) {
		t.Fatalf("expected UseColors() = true; got = %v", v)
	}
}

func Test_useColors_disableColors(t *testing.T) {
	defer withTestVar("NO_COLOR", "1")()

	if v := UseColors(); !cmp.Equal(v, &fa) {
		t.Fatalf("expected UseColors() = false; got = %v", v)
	}
}

func Test_useColors_default(t *testing.T) {
	defer withTestVar("NO_COLOR", "")()
	defer withTestVar("_KUBECTX_FORCE_COLOR", "")()

	if v := UseColors(); v != nil {
		t.Fatalf("expected UseColors() = nil; got=%v", *v)
	}
}
