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

func Error(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, ErrorColor.Sprint("error: ")+format+"\n", args...)
	return err
}

func Warning(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, WarningColor.Sprint("warning: ")+format+"\n", args...)
	return err
}

func Success(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, SuccessColor.Sprint("✔ ")+fmt.Sprintf(format+"\n", args...))
	return err
}
