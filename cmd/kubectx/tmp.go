// Copyright 2024 Google LLC
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

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/ahmetb/kubectx/internal/env"
	"github.com/ahmetb/kubectx/internal/kubeconfig"
)

// TmpOp wraps another operation with temporary kubeconfig mode.
type TmpOp struct {
	Inner Op
	Value string
}

func (op TmpOp) Run(stdout, stderr io.Writer) error {
	value := op.Value
	if value == "" {
		value = "1"
	}
	_ = os.Setenv(env.EnvTmp, value)
	if tmpPath, ok, err := kubeconfig.TempKubeconfigPath(); err == nil && ok {
		if kc := os.Getenv("KUBECONFIG"); kc != tmpPath {
			fmt.Fprint(stderr, tmpKubeconfigWarning(kc, tmpPath))
		}
	}
	return op.Inner.Run(stdout, stderr)
}

func tmpKubeconfigWarning(kubeconfigValue, tmpPath string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("warning: KUBECTX_TMP is set but KUBECONFIG is %q; to make kubectl follow the same temp context, run:\n  PowerShell:\n    $env:KUBECTX_TMP=%q\n    $env:KUBECONFIG=%q\n  cmd.exe:\n    set KUBECTX_TMP=%q\n    set KUBECONFIG=%q\n", kubeconfigValue, tmpPath, tmpPath, tmpPath, tmpPath)
	}
	return fmt.Sprintf("warning: KUBECTX_TMP is set but KUBECONFIG is %q; to make kubectl follow the same temp context, run:\n  export KUBECTX_TMP=%q\n  export KUBECONFIG=%q\n", kubeconfigValue, tmpPath, tmpPath)
}
