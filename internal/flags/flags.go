package flags

import (
	"fmt"
	"strings"
)

// Flags contains parsed flags.
type Flags struct {
	Delete           string
	Current          bool
	ShowHelp         bool
	Version          bool
	Unset            bool
	History          bool
	SelectedContext  []string
	NewContext       []string
}

// parseArgs parses given command line arguments to flags.
func parseArgs(args []string) (Flags, error) {
	var f Flags

	// filter out first argument (program name)
	if len(args) > 0 {
		args = args[1:]
	}

	// parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--help" || args[i] == "-h" {
			f.ShowHelp = true
			continue
		}

		if args[i] == "--current" || args[i] == "-c" {
			f.Current = true
			continue
		}

		if args[i] == "--unset" || args[i] == "-u" {
			f.Unset = true
			continue
		}

		if args[i] == "--version" || args[i] == "-V" {
			f.Version = true
			continue
		}

		if args[i] == "--history" {
			f.History = true
			continue
		}

		if (args[i] == "--delete" || args[i] == "-d") && i+1 < len(args) {
			f.Delete = args[i+1]
			i++
			continue
		}

		// <NEW>=<OLD>
		if strings.Contains(args[i], "=") {
			a := strings.SplitN(args[i], "=", 2)
			if len(a) != 2 {
				return f, fmt.Errorf("invalid argument: %s", args[i])
			}
			new, old := a[0], a[1]
			if new == "" || old == "" {
				return f, fmt.Errorf("invalid argument: %s", args[i])
			}
			f.NewContext = []string{new, old}
			continue
		}

		// <CONTEXT>
		f.SelectedContext = append(f.SelectedContext, args[i])
	}
	return f, nil
}

// New returns empty flags.
func New() *Flags {
	return &Flags{}
}

// Parse command line flags.
func (f *Flags) Parse(args []string) error {
	fl, err := parseArgs(args)
	if err != nil {
		return err
	}
	*f = fl
	return nil
}
