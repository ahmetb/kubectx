package main

import (
	"fmt"
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
		{name: "list shorthand",
			args: []string{"-l"},
			want: ListOp{}},
		{name: "list long form",
			args: []string{"--list"},
			want: ListOp{}},
		{name: "help shorthand",
			args: []string{"-h"},
			want: HelpOp{}},
		{name: "help long form",
			args: []string{"--help"},
			want: HelpOp{}},
		{name: "current shorthand",
			args: []string{"-c"},
			want: CurrentOp{}},
		{name: "current long form",
			args: []string{"--current"},
			want: CurrentOp{}},
		{name: "unset shorthand",
			args: []string{"-u"},
			want: UnsetOp{}},
		{name: "unset long form",
			args: []string{"--unset"},
			want: UnsetOp{}},
		{name: "switch by name",
			args: []string{"foo"},
			want: SwitchOp{Target: "foo"}},
		{name: "switch by swap",
			args: []string{"-"},
			want: SwitchOp{Target: "-"}},
		{name: "delete - without contexts",
			args: []string{"-d"},
			want: UnsupportedOp{fmt.Errorf("'-d' needs arguments")}},
		{name: "delete - current context",
			args: []string{"-d", "."},
			want: DeleteOp{[]string{"."}}},
		{name: "delete - multiple contexts",
			args: []string{"-d", ".", "a", "b"},
			want: DeleteOp{[]string{".", "a", "b"}}},
		{name: "rename context",
			args: []string{"a=b"},
			want: RenameOp{"a", "b"}},
		{name: "rename context with old=current",
			args: []string{"a=."},
			want: RenameOp{"a", "."}},
		{name: "unrecognized flag",
			args: []string{"-x"},
			want: UnsupportedOp{Err: fmt.Errorf("unsupported option '-x'")}},
		{name: "too many args",
			args: []string{"a", "b", "c"},
			want: UnsupportedOp{Err: fmt.Errorf("too many arguments")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseArgs(tt.args)

			var opts cmp.Options
			if _, ok := tt.want.(UnsupportedOp); ok {
				opts = append(opts, cmp.Comparer(func(x, y UnsupportedOp) bool {
					return (x.Err == nil && y.Err == nil) || (x.Err.Error() == y.Err.Error())
				}))
			}

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("parseArgs(%#v) diff: %s", tt.args, diff)
			}
		})
	}
}
