package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	var buf bytes.Buffer
	printHelp(&buf)

	out := buf.String()
	if !strings.Contains(out, "USAGE:") {
		t.Errorf("help string doesn't contain USAGE: ; output=%q", out)
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("does not end with new line; output=%q", out)
	}
}
