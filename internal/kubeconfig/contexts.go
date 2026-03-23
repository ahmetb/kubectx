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
	"slices"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func contextsNodeOf(fe *fileEntry) (*yaml.RNode, error) {
	contexts, err := fe.config.Pipe(yaml.Get("contexts"))
	if err != nil {
		return nil, err
	}
	if contexts == nil {
		return nil, errors.New("\"contexts\" entry is nil")
	} else if contexts.YNode().Kind != yaml.SequenceNode {
		return nil, errors.New("\"contexts\" is not a sequence node")
	}
	return contexts, nil
}

// contextNodeWithFileIndex searches for a context by name across all files.
// Returns the context node and the index of the file that contains it.
// Files without a valid "contexts" sequence are skipped, but if errors occur
// during lookup they are included in the final error message.
func (k *Kubeconfig) contextNodeWithFileIndex(name string) (*yaml.RNode, int, error) {
	var fileErrors []error
	for i := range k.files {
		contexts, err := contextsNodeOf(&k.files[i])
		if err != nil {
			fileErrors = append(fileErrors, fmt.Errorf("file %d: %w", i, err))
			continue
		}
		context, err := contexts.Pipe(yaml.Lookup("[name=" + name + "]"))
		if err != nil {
			fileErrors = append(fileErrors, fmt.Errorf("file %d lookup: %w", i, err))
			continue
		}
		if context != nil {
			return context, i, nil
		}
	}
	if len(fileErrors) > 0 {
		return nil, -1, fmt.Errorf("context with name %q not found (errors in files: %w)",
			name, errors.Join(fileErrors...))
	}
	return nil, -1, fmt.Errorf("context with name %q not found", name)
}

func (k *Kubeconfig) contextNode(name string) (*yaml.RNode, error) {
	node, _, err := k.contextNodeWithFileIndex(name)
	return node, err
}

func (k *Kubeconfig) ContextNames() ([]string, error) {
	seen := make(map[string]bool)
	var names []string

	for i := range k.files {
		contexts, err := k.files[i].config.Pipe(yaml.Get("contexts"))
		if err != nil {
			return nil, fmt.Errorf("failed to get contexts: %w", err)
		}
		if contexts == nil {
			continue
		}
		fileNames, err := contexts.ElementValues("name")
		if err != nil {
			return nil, fmt.Errorf("failed to get context names: %w", err)
		}
		for _, n := range fileNames {
			if !seen[n] {
				seen[n] = true
				names = append(names, n)
			}
		}
	}
	return names, nil
}

func (k *Kubeconfig) ContextExists(name string) (bool, error) {
	names, err := k.ContextNames()
	if err != nil {
		return false, err
	}
	return slices.Contains(names, name), nil
}
