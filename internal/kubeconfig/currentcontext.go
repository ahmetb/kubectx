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
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// GetCurrentContext returns "current-context" value from the first file
// that has a non-empty current-context, or returns ("", nil) if not found.
func (k *Kubeconfig) GetCurrentContext() (string, error) {
	for _, fe := range k.files {
		v, err := fe.config.Pipe(yaml.Get("current-context"))
		if err != nil {
			return "", fmt.Errorf("failed to read current-context: %w", err)
		}
		if s := yaml.GetValue(v); s != "" {
			return s, nil
		}
	}
	return "", nil
}

// UnsetCurrentContext clears the current-context field in the first file.
func (k *Kubeconfig) UnsetCurrentContext() error {
	if len(k.files) == 0 {
		return errNoFiles
	}
	return k.files[0].config.PipeE(yaml.SetField("current-context", yaml.NewStringRNode("")))
}
