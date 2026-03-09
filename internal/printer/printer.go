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
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	ErrorColor   = color.New(color.FgRed, color.Bold)
	WarningColor = color.New(color.FgYellow, color.Bold)
	SuccessColor = color.New(color.FgGreen)
)

func init() {
	colors := useColors()
	if colors == nil {
		return
	}
	if *colors {
		ErrorColor.EnableColor()
		WarningColor.EnableColor()
		SuccessColor.EnableColor()
	} else {
		ErrorColor.DisableColor()
		WarningColor.DisableColor()
		SuccessColor.DisableColor()
	}
}

func Error(w io.Writer, format string, args ...any) error {
	_, err := io.WriteString(w, ErrorColor.Sprint("error: ")+fmt.Sprintf(format, args...)+"\n")
	return err
}

func Warning(w io.Writer, format string, args ...any) error {
	_, err := io.WriteString(w, WarningColor.Sprint("warning: ")+fmt.Sprintf(format, args...)+"\n")
	return err
}

func Success(w io.Writer, format string, args ...any) error {
	_, err := io.WriteString(w, SuccessColor.Sprint("✔ ")+fmt.Sprintf(format, args...)+"\n")
	return err
}
