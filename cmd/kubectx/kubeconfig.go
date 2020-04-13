package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

var (
	defaultLoader kubeconfig.Loader = new(StandardKubeconfigLoader)
)

type StandardKubeconfigLoader struct{}

type kubeconfigFile struct{ *os.File }

func (kf *kubeconfigFile) Reset() error {
	if err := kf.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate file")
	}
	_, err := kf.Seek(0, 0)
	return errors.Wrap(err, "failed to seek in file")
}

func (*StandardKubeconfigLoader) Load() (kubeconfig.ReadWriteResetCloser, error) {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Errorf("kubeconfig file not found at %s", cfgPath)
		}
		return nil, errors.Wrap(err, "failed to open file")
	}
	return &kubeconfigFile{f}, nil
}

func kubeconfigPath() (string, error) {
	// KUBECONFIG env var
	if v := os.Getenv("KUBECONFIG"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO KUBECONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in KUBECONFIG are currently not supported")
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
