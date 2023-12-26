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
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// GetCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func (k *Kubeconfig) GetCurrentContext() string {
	v, err := k.config.Pipe(yaml.Get("current-context"))
	if err != nil {
		return ""
	}
	return yaml.GetValue(v)
}

func (k *Kubeconfig) UnsetCurrentContext() error {
	return k.config.PipeE(yaml.SetField("current-context", yaml.NewStringRNode("")))
}
