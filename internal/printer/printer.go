package printer

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	successColor = color.New(color.FgGreen)
)

func init() {
	colors := UseColors()
	if colors == nil {
		return
	}
	if *colors {
		errorColor.EnableColor()
		warningColor.EnableColor()
		successColor.EnableColor()
	} else {
		errorColor.DisableColor()
		warningColor.DisableColor()
		successColor.DisableColor()
	}
}

func Error(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.RedString("error: ")+format+"\n", args...)
}

func Warning(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.YellowString("warning: ")+format+"\n", args...)
}

func Success(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, color.GreenString(fmt.Sprintf(format+"\n", args...)))
}
