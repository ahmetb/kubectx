package kubeconfig

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ReadWriteResetCloser interface {
	io.ReadWriteCloser

	// Reset truncates the file and seeks to the beginning of the file.
	Reset() error
}

type Loader interface {
	Load(string) (ReadWriteResetCloser, error)
}

type Kubeconfig struct {
	loader Loader

	f        ReadWriteResetCloser
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

func (k *Kubeconfig) Parse() error {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return errors.Wrap(err, "cannot determine kubeconfig path")
	}

	f, err := k.loader.Load(cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to load")
	}

	k.f = f
	var v yaml.Node
	if err := yaml.NewDecoder(f).Decode(&v); err != nil {
		return errors.Wrap(err, "failed to decode")
	}
	k.rootNode = v.Content[0]
	if k.rootNode.Kind != yaml.MappingNode {
		return errors.New("kubeconfig file is not a map document")
	}
	return nil
}

func (k *Kubeconfig) Bytes() ([]byte, error) {
	return yaml.Marshal(k.rootNode)
}

func (k *Kubeconfig) Save() error {
	if err := k.f.Reset(); err != nil {
		return errors.Wrap(err, "failed to reset file")
	}
	enc := yaml.NewEncoder(k.f)
	enc.SetIndent(0)
	return enc.Encode(k.rootNode)
}
