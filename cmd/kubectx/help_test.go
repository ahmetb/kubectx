package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer
	if err := (&HelpOp{}).Run(&buf, &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "USAGE:") {
		t.Errorf("help string doesn't contain USAGE: ; output=\"%s\"", out)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("does not end with New line; output=\"%s\"", out)
	}
}
