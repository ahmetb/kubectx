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

package kubeconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

var (
	DefaultLoader Loader = new(StandardKubeconfigLoader)
)

type StandardKubeconfigLoader struct{}

type kubeconfigFile struct{ *os.File }

func (*StandardKubeconfigLoader) Load() ([]ReadWriteResetCloser, error) {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return nil, fmt.Errorf("cannot determine kubeconfig path: %w", err)
	}

	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// TODO we'll return all kubeconfig files when we start implementing multiple kubeconfig support
	return []ReadWriteResetCloser{ReadWriteResetCloser(&kubeconfigFile{f})}, nil
}

func (kf *kubeconfigFile) Reset() error {
	if err := kf.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := kf.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek in file: %w", err)
	}
	return nil
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
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "config"), nil
}
