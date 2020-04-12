package main

import (
	"os"
	"testing"
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

func Test_useColors_forceColors(t *testing.T) {
	defer withTestVar("_KUBECTX_FORCE_COLOR", "1")()
	defer withTestVar("NO_COLOR", "1")()

	if !useColors() {
		t.Fatal("expected useColors() = true")
	}
}

func Test_useColors_disableColors(t *testing.T) {
	defer withTestVar("NO_COLOR", "1")()

	if useColors() {
		t.Fatal("expected useColors() = false")
	}
}

func Test_useColors_default(t *testing.T) {
	defer withTestVar("NO_COLOR", "")()
	defer withTestVar("_KUBECTX_FORCE_COLOR", "")()

	if !useColors() {
		t.Fatal("expected useColors() = true")
	}
}
