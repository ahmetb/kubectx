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

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (k *Kubeconfig) DeleteContextEntry(deleteName string) error {
	_, fileIdx, err := k.contextNodeWithFileIndex(deleteName)
	if err != nil {
		return err
	}

	contexts, err := contextsNodeOf(&k.files[fileIdx])
	if err != nil {
		return err
	}
	return contexts.PipeE(
		yaml.ElementSetter{
			Keys:   []string{"name"},
			Values: []string{deleteName},
		},
	)
}

// ModifyCurrentContext always writes to the first file (matching kubectl behavior).
func (k *Kubeconfig) ModifyCurrentContext(name string) error {
	if len(k.files) == 0 {
		return errNoFiles
	}
	return k.files[0].config.PipeE(yaml.SetField("current-context", yaml.NewScalarRNode(name)))
}

func (k *Kubeconfig) ModifyContextName(old, new string) error {
	context, _, err := k.contextNodeWithFileIndex(old)
	if err != nil {
		return err
	}
	if context == nil {
		return errors.New("\"contexts\" entry is nil")
	}
	return context.PipeE(yaml.SetField("name", yaml.NewScalarRNode(new)))
}
