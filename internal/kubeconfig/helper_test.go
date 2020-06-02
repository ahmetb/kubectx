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
