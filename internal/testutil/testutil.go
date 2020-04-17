package testutil

import "os"

// WithEnvVar sets an env var temporarily. Call its return value
// in defer to restore original value in env (if exists).
func WithEnvVar(key, value string) func() {
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

