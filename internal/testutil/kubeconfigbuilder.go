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

package testutil

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type ContextObj struct {
	Name    string `yaml:"name,omitempty"`
	Context struct {
		Namespace string `yaml:"namespace,omitempty"`
		User      string `yaml:"user,omitempty"`
		Cluster   string `yaml:"cluster,omitempty"`
	} `yaml:"context,omitempty"`
}

func Ctx(name string) *ContextObj                        { return &ContextObj{Name: name} }
func (c *ContextObj) Ns(ns string) *ContextObj           { c.Context.Namespace = ns; return c }
func (c *ContextObj) User(user string) *ContextObj       { c.Context.User = user; return c }
func (c *ContextObj) Cluster(cluster string) *ContextObj { c.Context.Cluster = cluster; return c }

type UserObj struct {
	Name string `yaml:"name,omitempty"`
}

func User(name string) *UserObj { return &UserObj{Name: name} }

type ClusterObj struct {
	Name string `yaml:"name,omitempty"`
}

func Cluster(name string) *ClusterObj { return &ClusterObj{Name: name} }

type Kubeconfig map[string]interface{}

func KC() *Kubeconfig {
	return &Kubeconfig{
		"apiVersion": "v1",
		"kind":       "Config"}
}

func (k *Kubeconfig) Set(key string, v interface{}) *Kubeconfig { (*k)[key] = v; return k }
func (k *Kubeconfig) WithCurrentCtx(s string) *Kubeconfig       { (*k)["current-context"] = s; return k }
func (k *Kubeconfig) WithCtxs(c ...*ContextObj) *Kubeconfig     { (*k)["contexts"] = c; return k }
func (k *Kubeconfig) WithUsers(u ...*UserObj) *Kubeconfig       { (*k)["users"] = u; return k }
func (k *Kubeconfig) WithClusters(c ...*ClusterObj) *Kubeconfig { (*k)["clusters"] = c; return k }

func (k *Kubeconfig) ToYAML(t *testing.T) string {
	t.Helper()
	var v strings.Builder
	if err := yaml.NewEncoder(&v).Encode(*k); err != nil {
		t.Fatalf("failed to encode mock kubeconfig: %v", err)
	}
	return v.String()
}
