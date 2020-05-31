package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ahmetb/kubectx/internal/cmdutil"
)

var defaultDir = filepath.Join(cmdutil.HomeDir(), ".kube", "kubens")

type NSFile struct {
	dir string
	ctx string
}

func NewNSFile(ctx string) NSFile { return NSFile{dir: defaultDir, ctx: ctx} }

func (f NSFile) path() string {
	fn := f.ctx
	if isWindows() {
		// bug 230: eks clusters contain ':' in ctx name, not a valid file name for win32
		fn = strings.ReplaceAll(fn, ":", "__")
	}
	return filepath.Join(f.dir, fn)
}

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

// isWindows determines if the process is running on windows OS.
func isWindows() bool {
	if os.Getenv("_FORCE_GOOS") == "windows" { // for testing
		return true
	}
	return runtime.GOOS == "windows"
}
