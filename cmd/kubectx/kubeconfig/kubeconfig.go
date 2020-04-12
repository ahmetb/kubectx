package kubeconfig

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ReadWriteResetCloser interface{
	io.ReadWriteCloser

	// Reset truncates the file and seeks to the beginning of the file.
	Reset() error
}

type Loader interface {
	Load() (ReadWriteResetCloser, error)
}

type Kubeconfig struct {
	loader   Loader

	f ReadWriteResetCloser
	rootNode *yaml.Node
}

func (k *Kubeconfig) WithLoader(l Loader) *Kubeconfig {
	k.loader = l
	return k
}

func (k *Kubeconfig) Close() error {
	if k.f == nil {
		return nil
	}
	return k.f.Close()
}

func (k *Kubeconfig) ParseRaw() (*yaml.Node, error) {
	f, err := k.loader.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubeconfig")
	}

	k.f = f

	var v yaml.Node
	if err := yaml.NewDecoder(f).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "failed to decode kubeconfig")
	}
	k.rootNode = v.Content[0]
	return k. rootNode, nil
}

func (k *Kubeconfig) Save() error {
	if err := k.f.Reset(); err != nil {
		return errors.Wrap(err, "failed to reset file")
	}
	enc := yaml.NewEncoder(k.f)
	enc.SetIndent(2)
	return enc.Encode(k.rootNode)
}
