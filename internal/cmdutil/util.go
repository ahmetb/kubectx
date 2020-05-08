package cmdutil

import (
	"os"

	"github.com/pkg/errors"
)

func HomeDir() string {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return v
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}
	return home
}

// IsNotFoundErr determines if the underlying error is os.IsNotExist. Right now
// errors from github.com/pkg/errors doesn't work with os.IsNotExist.
func IsNotFoundErr(err error) bool {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if os.IsNotExist(e) {
			return true
		}
	}
	return false
}
