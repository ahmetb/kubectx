package cmdutil

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintDeprecatedEnvWarnings_noDeprecatedVars(t *testing.T) {
	var out bytes.Buffer
	PrintDeprecatedEnvWarnings(&out, []string{
		"A=B",
		"PATH=/foo:/bar:/bin",
	})
	if v := out.String(); len(v) > 0 {
		t.Fatalf("something written to buf: %v", v)
	}
}

func TestPrintDeprecatedEnvWarnings_bgColors(t *testing.T) {
	var out bytes.Buffer

	PrintDeprecatedEnvWarnings(&out, []string{
		"KUBECTX_CURRENT_FGCOLOR=1",
		"KUBECTX_CURRENT_BGCOLOR=2",
	})
	v := out.String()
	if !strings.Contains(v, "KUBECTX_CURRENT_FGCOLOR") {
		t.Fatalf("output doesn't contain 'KUBECTX_CURRENT_FGCOLOR': \"%s\"", v)
	}
	if !strings.Contains(v, "KUBECTX_CURRENT_BGCOLOR") {
		t.Fatalf("output doesn't contain 'KUBECTX_CURRENT_BGCOLOR': \"%s\"", v)
	}
}
