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
