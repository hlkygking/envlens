package envdiff

import (
	"fmt"
	"strings"
)

// FormatOptions controls how diff output is rendered.
type FormatOptions struct {
	ShowUnchanged bool
	Colorize      bool
	Compact       bool
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		ShowUnchanged: false,
		Colorize:      false,
		Compact:       false,
	}
}

// FormatEntry renders a single DiffEntry as a human-readable line.
func FormatEntry(e Entry, opts FormatOptions) string {
	switch e.Status {
	case StatusAdded:
		line := fmt.Sprintf("+ %s=%s", e.Key, e.NewValue)
		if opts.Colorize {
			return colorGreen(line)
		}
		return line
	case StatusRemoved:
		line := fmt.Sprintf("- %s=%s", e.Key, e.OldValue)
		if opts.Colorize {
			return colorRed(line)
		}
		return line
	case StatusModified:
		if opts.Compact {
			line := fmt.Sprintf("~ %s=%s", e.Key, e.NewValue)
			if opts.Colorize {
				return colorYellow(line)
			}
			return line
		}
		line := fmt.Sprintf("~ %s: %s -> %s", e.Key, e.OldValue, e.NewValue)
		if opts.Colorize {
			return colorYellow(line)
		}
		return line
	case StatusUnchanged:
		if !opts.ShowUnchanged {
			return ""
		}
		return fmt.Sprintf("  %s=%s", e.Key, e.NewValue)
	}
	return ""
}

// FormatAll renders all entries using the given options, skipping blank lines.
func FormatAll(entries []Entry, opts FormatOptions) string {
	var lines []string
	for _, e := range entries {
		if line := FormatEntry(e, opts); line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return "(no changes)"
	}
	return strings.Join(lines, "\n")
}

func colorGreen(s string) string  { return "\033[32m" + s + "\033[0m" }
func colorRed(s string) string    { return "\033[31m" + s + "\033[0m" }
func colorYellow(s string) string { return "\033[33m" + s + "\033[0m" }
