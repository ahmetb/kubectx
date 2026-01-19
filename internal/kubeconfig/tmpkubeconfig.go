// Copyright 2024 Google LLC
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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/pkg/errors"
)

// TempKubeconfigPath returns the temp kubeconfig path if configured.
func TempKubeconfigPath() (string, bool, error) {
	raw := strings.TrimSpace(os.Getenv(env.EnvTmp))
	if raw == "" {
		return "", false, nil
	}

	switch strings.ToLower(raw) {
	case "0", "false", "no", "off":
		return "", false, nil
	case "1", "true", "yes", "on", "auto":
		dir := os.Getenv("XDG_RUNTIME_DIR")
		if dir == "" {
			dir = os.TempDir()
		}
		if dir == "" {
			return "", false, errors.New("temporary directory not available")
		}
		name := fmt.Sprintf("kubectx-%d", os.Getppid())
		return filepath.Join(dir, "kubectx", name), true, nil
	default:
		return raw, true, nil
	}
}

func ensureTmpKubeconfig(basePath, tmpPath string) error {
	if _, err := os.Stat(tmpPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to stat temp kubeconfig")
	}

	if err := os.MkdirAll(filepath.Dir(tmpPath), 0700); err != nil {
		return errors.Wrap(err, "failed to create temp kubeconfig directory")
	}

	src, err := os.Open(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, "kubeconfig file not found")
		}
		return errors.Wrap(err, "failed to open kubeconfig")
	}
	defer src.Close()

	dst, err := os.OpenFile(tmpPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to create temp kubeconfig")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrap(err, "failed to copy kubeconfig")
	}
	return nil
}
