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

package printer

import (
	"os"

	"github.com/fatih/color"

	"github.com/ahmetb/kubectx/internal/env"
)

var (
	ActiveItemColor = color.New(color.FgGreen, color.Bold)
)

func init() {
	EnableOrDisableColor(ActiveItemColor)
}

// useColors returns true if colors are force-enabled,
// false if colors are disabled, or nil for default behavior
// which is determined based on factors like if stdout is tty.
func useColors() *bool {
	tr, fa := true, false
	if os.Getenv(env.EnvForceColor) != "" {
		return &tr
	} else if os.Getenv(env.EnvNoColor) != "" {
		return &fa
	}
	return nil
}

// EnableOrDisableColor determines if color should be force-enabled or force-disabled
// or left untouched based on environment configuration.
func EnableOrDisableColor(c *color.Color) {
	if v := useColors(); v != nil && *v {
		c.EnableColor()
	} else if v != nil && !*v {
		c.DisableColor()
	}
}
