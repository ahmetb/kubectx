package cmdutil

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

var (
	DefaultLoader kubeconfig.Loader = new(StandardKubeconfigLoader)
)

type StandardKubeconfigLoader struct{}

type kubeconfigFile struct{ *os.File }

func (*StandardKubeconfigLoader) Load() (kubeconfig.ReadWriteResetCloser, error) {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, "kubeconfig file not found")
		}
		return nil, errors.Wrap(err, "failed to open file")
	}
	return &kubeconfigFile{f}, nil
}

func (kf *kubeconfigFile) Reset() error {
	if err := kf.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate file")
	}
	_, err := kf.Seek(0, 0)
	return errors.Wrap(err, "failed to seek in file")
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

	// default path
	home := HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "config"), nil
}

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
