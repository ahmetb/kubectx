package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseArgs_new(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Op
	}{
		{name: "nil Args",
			args: nil,
			want: ListOp{}},
		{name: "empty Args",
			args: []string{},
			want: ListOp{}},
		{name: "help shorthand",
			args: []string{"-h"},
			want: HelpOp{}},
		{name: "help long form",
			args: []string{"--help"},
			want: HelpOp{}},
		{name: "switch by name",
			args: []string{"foo"},
			want: SwitchOp{Target: "foo"}},
		{name: "switch by swap",
			args: []string{"-"},
			want: SwitchOp{Target: "-"}},
		{name: "unrecognized flag",
			args: []string{"-x"},
			want: UnknownOp{Args: []string{"-x"}}},
		// TODO add more UnknownOp cases

		// TODO consider these cases
		// - kubectx foo --help
		// - kubectx -h --help
		// - kubectx -d foo --h
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseArgs(tt.args)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("parseArgs(%#v) diff: %s", tt.args, diff)
			}
		})
	}
}
