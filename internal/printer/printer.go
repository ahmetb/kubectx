package printer

import (
	"fmt"
	"io"
	"strings"

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

func Error(w io.Writer, format string, args ...interface{}) error {
	return writeOutput(w, ErrorColor.Sprint("error: ")+format+"\n", args...)
}

func Warning(w io.Writer, format string, args ...interface{}) error {
	return writeOutput(w, WarningColor.Sprint("warning: ")+format+"\n", args...)
}

func Success(w io.Writer, format string, args ...interface{}) error {
	return writeOutput(w, SuccessColor.Sprint("✔ ")+format+"\n", args...)
}

func writeOutput(w io.Writer, format string, args ...interface{}) error {
	// Replace %q with "%s" so unescaped color sequences are written to output
	format = strings.ReplaceAll(format, "%q", "\"%s\"")

	_, err := fmt.Fprintf(w, format, args...)
	return err
}
