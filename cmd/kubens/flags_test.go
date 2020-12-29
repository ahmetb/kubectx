// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		{name: "switch by name",
			args: []string{"foo"},
			want: SwitchOp{Target: "foo"}},
		{name: "switch by swap",
			args: []string{"-"},
			want: SwitchOp{Target: "-"}},
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
