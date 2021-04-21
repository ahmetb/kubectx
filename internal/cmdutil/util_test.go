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

package cmdutil

import (
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func Test_homeDir(t *testing.T) {
	type env struct{ k, v string }
	cases := []struct {
		name string
		envs []env
		want string
	}{
		{
			name: "XDG_CACHE_HOME precedence",
			envs: []env{
				{"XDG_CACHE_HOME", "xdg"},
				{"HOME", "home"},
			},
			want: "xdg",
		},
		{
			name: "HOME over USERPROFILE",
			envs: []env{
				{"HOME", "home"},
				{"USERPROFILE", "up"},
			},
			want: "home",
		},
		{
			name: "only USERPROFILE available",
			envs: []env{
				{"XDG_CACHE_HOME", ""},
				{"HOME", ""},
				{"USERPROFILE", "up"},
			},
			want: "up",
		},
		{
			name: "none available",
			envs: []env{
				{"XDG_CACHE_HOME", ""},
				{"HOME", ""},
				{"USERPROFILE", ""},
			},
			want: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			var unsets []func()
			for _, e := range c.envs {
				unsets = append(unsets, testutil.WithEnvVar(e.k, e.v))
			}

			got := HomeDir()
			if got != c.want {
				t.Errorf("expected:%q got:%q", c.want, got)
			}
			for _, u := range unsets {
				u()
			}
		})
	}
}
