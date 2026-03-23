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
	"io"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ReadWriteResetCloser interface {
	io.ReadWriteCloser

	// Reset truncates the file and seeks to the beginning of the file.
	Reset() error
}

type Loader interface {
	Load() ([]ReadWriteResetCloser, error)
}

type fileEntry struct {
	f      ReadWriteResetCloser
	config *yaml.RNode
}

type Kubeconfig struct {
	loader Loader
	files  []fileEntry
}

var errNoFiles = errors.New("no kubeconfig files loaded")

func (k *Kubeconfig) WithLoader(l Loader) *Kubeconfig {
	k.loader = l
	return k
}

func (k *Kubeconfig) Close() error {
	var firstErr error
	for _, fe := range k.files {
		if err := fe.f.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (k *Kubeconfig) Parse() error {
	rwcs, err := k.loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	k.files = make([]fileEntry, 0, len(rwcs))
	for i, f := range rwcs {
		var v yaml.Node
		if err := yaml.NewDecoder(f).Decode(&v); err != nil {
			// Close all file handles on failure to avoid leaks.
			for _, rf := range rwcs {
				rf.Close()
			}
			return fmt.Errorf("failed to decode file %d: %w", i, err)
		}
		rn := yaml.NewRNode(&v)
		if rn.YNode().Kind != yaml.MappingNode {
			for _, rf := range rwcs {
				rf.Close()
			}
			return fmt.Errorf("kubeconfig file %d is not a map document", i)
		}
		k.files = append(k.files, fileEntry{f: f, config: rn})
	}
	return nil
}

func (k *Kubeconfig) Bytes() ([]byte, error) {
	if len(k.files) == 0 {
		return nil, errNoFiles
	}

	if len(k.files) == 1 {
		str, err := k.files[0].config.String()
		if err != nil {
			return nil, err
		}
		return []byte(str), nil
	}

	// Build a merged config for multi-file case.
	// Start with a copy of the first file's structure.
	merged := k.files[0].config.Copy()

	// Merge contexts, clusters, and users from all files (first wins for duplicates).
	for _, key := range []string{"contexts", "clusters", "users"} {
		mergedSeq, err := mergeSequences(k.files, key)
		if err != nil {
			return nil, fmt.Errorf("failed to merge %s: %w", key, err)
		}
		if mergedSeq != nil {
			if err := merged.PipeE(yaml.SetField(key, mergedSeq)); err != nil {
				return nil, err
			}
		}
	}

	// Use the first non-empty current-context.
	cur, err := k.GetCurrentContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context for merge: %w", err)
	}
	if cur != "" {
		if err := merged.PipeE(yaml.SetField("current-context", yaml.NewScalarRNode(cur))); err != nil {
			return nil, err
		}
	}

	str, err := merged.String()
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// mergeSequences merges a named sequence field (e.g. "contexts") across multiple files.
// The first occurrence of each entry (by "name" key) wins.
// Files where the field is missing or not a sequence are silently skipped (matching kubectl merge behavior).
func mergeSequences(files []fileEntry, field string) (*yaml.RNode, error) {
	seen := make(map[string]bool)
	var elements []*yaml.RNode

	for _, fe := range files {
		seq, err := fe.config.Pipe(yaml.Get(field))
		if err != nil || seq == nil {
			continue
		}
		if seq.YNode().Kind != yaml.SequenceNode {
			continue
		}
		for _, elem := range seq.YNode().Content {
			rn := yaml.NewRNode(elem)
			name, err := rn.Pipe(yaml.Get("name"))
			if err != nil || name == nil {
				continue
			}
			n := yaml.GetValue(name)
			if n != "" && seen[n] {
				continue
			}
			seen[n] = true
			elements = append(elements, rn)
		}
	}

	if len(elements) == 0 {
		return nil, nil
	}
	seqNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, elem := range elements {
		seqNode.Content = append(seqNode.Content, elem.YNode())
	}
	return yaml.NewRNode(seqNode), nil
}

func (k *Kubeconfig) Save() error {
	for i := range k.files {
		if err := k.files[i].f.Reset(); err != nil {
			return fmt.Errorf("failed to reset file %d: %w", i, err)
		}
		enc := yaml.NewEncoder(k.files[i].f)
		enc.SetIndent(0)
		if err := enc.Encode(k.files[i].config.YNode()); err != nil {
			return fmt.Errorf("failed to encode file %d: %w", i, err)
		}
		if err := enc.Close(); err != nil {
			return fmt.Errorf("failed to close encoder for file %d: %w", i, err)
		}
	}
	return nil
}
