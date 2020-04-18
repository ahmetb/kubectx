package testutil

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type Context struct {
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

func Ctx(name string) *Context           { return &Context{Name: name} }
func (c *Context) Ns(ns string) *Context { c.Namespace = ns; return c }

type Kubeconfig map[string]interface{}

func KC() *Kubeconfig {
	return &Kubeconfig{
		"apiVersion": "v1",
		"kind":       "Config"}
}

func (k *Kubeconfig) Set(key string, v interface{}) *Kubeconfig { (*k)[key] = v; return k }
func (k *Kubeconfig) WithCurrentCtx(s string) *Kubeconfig       { (*k)["current-context"] = s; return k }
func (k *Kubeconfig) WithCtxs(c ...*Context) *Kubeconfig        { (*k)["contexts"] = c; return k }

func (k *Kubeconfig) ToYAML(t *testing.T) string {
	t.Helper()
	var v strings.Builder
	if err := yaml.NewEncoder(&v).Encode(*k); err != nil {
		t.Fatalf("failed to encode mock kubeconfig: %v", err)
	}
	return v.String()
}
