package cmdutil

import (
	"testing"

	"github.com/ahmetb/kubectx/internal/testutil"
)

func Test_fuzzyPicker(t *testing.T) {
	type env struct{ k, v string }

	cases := []struct {
		name string
		envs []env
		want string
	}{
		{
			name: "PICKER is fzf",
			envs: []env{
				{"PICKER", "fzf"},
			},
			want: "fzf",
		}, {
			name: "PICKER is sk",
			envs: []env{
				{"PICKER", "sk"},
			},
			want: "sk",
		}, {
			name: "PICKER is not set",
			envs: []env{},
			want: "fzf",
		}, {
			name: "PICKER is other than fzf and sk",
			envs: []env{{"PICKER", "other-fuzzer"}},
			want: "fzf",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			var unsets []func()
			for _, e := range c.envs {
				unsets = append(unsets, testutil.WithEnvVar(e.k, e.v))
			}
			got := fuzzyPicker()
			if got != c.want {
				t.Errorf("want: %s, got: %s", c.want, got)
			}
			for _, u := range unsets {
				u()
			}
		})
	}
}
