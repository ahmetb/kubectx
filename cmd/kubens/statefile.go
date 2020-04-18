package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

var defaultDir = filepath.Join(cmdutil.HomeDir(), ".kube", "kubens")

type NSFile struct {
	dir string
	ctx string
}

func NewNSFile(ctx string) NSFile { return NSFile{dir: defaultDir, ctx: ctx} }

func (f NSFile) path() string { return filepath.Join(f.dir, f.ctx) }

// Load reads the previous namespace setting, or returns empty if not exists.
func (f NSFile) Load() (string, error) {
	b, err := ioutil.ReadFile(f.path())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(bytes.TrimSpace(b)), nil
}

// Save stores the previous namespace information in the file.
func (f NSFile) Save(value string) error {
	d := filepath.Dir(f.path())
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(f.path(), []byte(value), 0644)
}
