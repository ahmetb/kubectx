package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func kubeconfigPath() (string, error) {
	// KUBECONFIG env var
	if v := os.Getenv("KUBECONFIG"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO KUBECONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in KUBECONFIG currently not supported")
		}
		return v, nil
	}

	home := homeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}

	// return default path
	return filepath.Join(home, ".kube", "config"), nil
}

func homeDir() string {
	// TODO move tests for this out of kubeconfigPath to TestHomeDir()
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return v
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}
	return home
}


// TODO parseKubeconfig doesn't seem necessary when there's a raw version that returns what's needed
func parseKubeconfig(path string) (kubeconfig, error) {
	// TODO refactor to accept io.Reader instead of file
	var v kubeconfig

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return v, nil
		}
		return v, errors.Wrap(err, "file open error")
	}
	err = yaml.NewDecoder(f).Decode(&v)
	return v, errors.Wrap(err, "yaml parse error")
}
