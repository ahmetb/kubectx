package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func parseKubeconfigRaw(r io.Reader) (*yaml.Node, error) {
	// TODO DELETE
	var v yaml.Node
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return v.Content[0], nil
}

func saveKubeconfigRaw(w io.Writer, rootNode *yaml.Node) error {
	if f, ok := w.(*os.File); ok {
		if err := f.Truncate(0); err != nil {
			return errors.Wrap(err, "failed to truncate")
		}
		if _, err := f.Seek(0, 0); err != nil {
			return errors.Wrap(err, "failed to seek")
		}
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(rootNode)
}

func openKubeconfig() (f *os.File, rootNode *yaml.Node, err error) {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}
	f, err = os.OpenFile(cfgPath, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open file")
	}

	kc, err := parseKubeconfigRaw(f)
	if err != nil {
		f.Close()
		return nil, nil, errors.Wrap(err, "yaml parse error")
	}
	return f, kc, nil
}
