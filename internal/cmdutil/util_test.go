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
	"path/filepath"
	"testing"
)

func Test_homeDir(t *testing.T) {
	type env struct{ k, v string }
	cases := []struct {
		name string
		envs []env
		want string
	}{
		{
			name: "don't use XDG_CACHE_HOME as homedir",
			envs: []env{
				{"XDG_CACHE_HOME", "xdg"},
				{"HOME", "home"},
			},
			want: "home",
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
				{"HOME", ""},
				{"USERPROFILE", "up"},
			},
			want: "up",
		},
		{
			name: "none available",
			envs: []env{
				{"HOME", ""},
				{"USERPROFILE", ""},
			},
			want: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			for _, e := range c.envs {
				tt.Setenv(e.k, e.v)
			}

			got := HomeDir()
			if got != c.want {
				t.Errorf("expected:%q got:%q", c.want, got)
			}
		})
	}
}

func TestCacheDir(t *testing.T) {
	t.Run("XDG_CACHE_HOME set", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "/tmp/xdg-cache")
		t.Setenv("HOME", "/home/user")
		if got := CacheDir(); got != "/tmp/xdg-cache" {
			t.Errorf("expected:%q got:%q", "/tmp/xdg-cache", got)
		}
	})
	t.Run("XDG_CACHE_HOME unset, falls back to HOME/.kube", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		t.Setenv("HOME", "/home/user")
		want := filepath.Join("/home/user", ".kube")
		if got := CacheDir(); got != want {
			t.Errorf("expected:%q got:%q", want, got)
		}
	})
	t.Run("neither set", func(t *testing.T) {
		t.Setenv("XDG_CACHE_HOME", "")
		t.Setenv("HOME", "")
		t.Setenv("USERPROFILE", "")
		if got := CacheDir(); got != "" {
			t.Errorf("expected:%q got:%q", "", got)
		}
	})
}
