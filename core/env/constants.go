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

package env

const (
	// EnvFZFIgnore describes the environment variable to set to disable
	// interactive context selection when fzf is installed.
	EnvFZFIgnore = "KUBECTX_IGNORE_FZF"

	// EnvForceColor describes the environment variable to disable color usage
	// when printing current context in a list.
	EnvNoColor = `NO_COLOR`

	// EnvForceColor describes the "internal" environment variable to force
	// color usage to show current context in a list.
	EnvForceColor = `_KUBECTX_FORCE_COLOR`

	// EnvDebug describes the internal environment variable for more verbose logging.
	EnvDebug = `DEBUG`
)
