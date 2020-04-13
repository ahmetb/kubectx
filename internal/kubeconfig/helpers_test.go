package kubeconfig

import (
	"bytes"
	"strings"
)

type testLoader struct {
	in  *strings.Reader
	out bytes.Buffer
}

func (t *testLoader) Read(p []byte) (n int, err error)    { return t.in.Read(p) }
func (t *testLoader) Write(p []byte) (n int, err error)   { return t.out.Write(p) }
func (t *testLoader) Close() error                        { return nil }
func (t *testLoader) Reset() error                        { return nil }
func (t *testLoader) Load() (ReadWriteResetCloser, error) { return t, nil }
