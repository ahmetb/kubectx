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
	"bytes"
	"io"
	"strings"
)

type MockKubeconfigLoader struct {
	in  io.Reader
	out bytes.Buffer
}

func (t *MockKubeconfigLoader) Read(p []byte) (n int, err error)  { return t.in.Read(p) }
func (t *MockKubeconfigLoader) Write(p []byte) (n int, err error) { return t.out.Write(p) }
func (t *MockKubeconfigLoader) Close() error                      { return nil }
func (t *MockKubeconfigLoader) Reset() error                      { return nil }
func (t *MockKubeconfigLoader) Load() ([]ReadWriteResetCloser, error) {
	return []ReadWriteResetCloser{ReadWriteResetCloser(t)}, nil
}
func (t *MockKubeconfigLoader) Output() string { return t.out.String() }

func WithMockKubeconfigLoader(kubecfg string) *MockKubeconfigLoader {
	return &MockKubeconfigLoader{in: strings.NewReader(kubecfg)}
}

// mockFile is a single in-memory kubeconfig file for multi-file testing.
type mockFile struct {
	in  io.Reader
	out bytes.Buffer
}

func (m *mockFile) Read(p []byte) (n int, err error)  { return m.in.Read(p) }
func (m *mockFile) Write(p []byte) (n int, err error) { return m.out.Write(p) }
func (m *mockFile) Close() error                      { return nil }
func (m *mockFile) Reset() error                      { return nil }

// MockMultiKubeconfigLoader implements Loader for testing with multiple kubeconfig files.
type MockMultiKubeconfigLoader struct {
	files []*mockFile
}

func (m *MockMultiKubeconfigLoader) Load() ([]ReadWriteResetCloser, error) {
	out := make([]ReadWriteResetCloser, len(m.files))
	for i, f := range m.files {
		out[i] = f
	}
	return out, nil
}

func (m *MockMultiKubeconfigLoader) OutputOf(index int) string {
	return m.files[index].out.String()
}

func WithMockMultiKubeconfigLoader(configs ...string) *MockMultiKubeconfigLoader {
	files := make([]*mockFile, len(configs))
	for i, c := range configs {
		files[i] = &mockFile{in: strings.NewReader(c)}
	}
	return &MockMultiKubeconfigLoader{files: files}
}
