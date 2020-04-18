package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

func kubectxPrevCtxFile() (string, error) {
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "kubectx"), nil
}

// readLastContext returns the saved previous context
// if the state file exists, otherwise returns "".
func readLastContext(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	return string(b), err
}

// writeLastContext saves the specified value to the state file.
// It creates missing parent directories.
func writeLastContext(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "failed to create parent directories")
	}
	return ioutil.WriteFile(path, []byte(value), 0644)
}
