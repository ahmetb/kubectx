package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseRenameSyntax(t *testing.T) {

	type out struct {
		New string
		Old string
		OK  bool
	}
	tests := []struct {
		name string
		in   string
		want out
	}{
		{
			name: "no equals sign",
			in:   "foo",
			want: out{OK: false},
		},
		{
			name: "no left side",
			in:   "=a",
			want: out{OK: false},
		},
		{
			name: "no right side",
			in:   "a=",
			want: out{OK: false},
		},
		{
			name: "correct format",
			in:   "a=b",
			want: out{
				New: "a",
				Old: "b",
				OK:  true,
			},
		},
		{
			name: "correct format with current context",
			in:   "NEW_NAME=.",
			want: out{
				New: "NEW_NAME",
				Old: ".",
				OK:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			new, old, ok := parseRenameSyntax(tt.in)
			got := out{
				New: new,
				Old: old,
				OK:  ok,
			}
			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf("parseRenameSyntax() diff=%s", diff)
			}
		})
	}
}
