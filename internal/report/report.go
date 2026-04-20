package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/diff"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextReport writes a human-readable diff report to w.
func TextReport(w io.Writer, entries []diff.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	added := filterByStatus(entries, diff.Added)
	removed := filterByStatus(entries, diff.Removed)
	modified := filterByStatus(entries, diff.Modified)

	if len(added) > 0 {
		fmt.Fprintln(w, "Added:")
		for _, e := range added {
			fmt.Fprintf(w, "  + %s=%s\n", e.Key, e.NewValue)
		}
	}

	if len(removed) > 0 {
		fmt.Fprintln(w, "Removed:")
		for _, e := range removed {
			fmt.Fprintf(w, "  - %s=%s\n", e.Key, e.OldValue)
		}
	}

	if len(modified) > 0 {
		fmt.Fprintln(w, "Modified:")
		for _, e := range modified {
			fmt.Fprintf(w, "  ~ %s: %q -> %q\n", e.Key, e.OldValue, e.NewValue)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d added, %d removed, %d modified\n",
		len(added), len(removed), len(modified))
}

// JSONReport writes a JSON diff report to w.
// If entries is empty, it writes an empty JSON array.
func JSONReport(w io.Writer, entries []diff.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "[]")
		return
	}

	fmt.Fprintln(w, "[")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "  {\"key\": %q, \"status\": %q, \"old\": %q, \"new\": %q}%s\n",
			e.Key, strings.ToLower(string(e.Status)), e.OldValue, e.NewValue, comma)
	}
	fmt.Fprintln(w, "]")
}

func filterByStatus(entries []diff.Entry, status diff.Status) []diff.Entry {
	var out []diff.Entry
	for _, e := range entries {
		if e.Status == status {
			out = append(out, e)
		}
	}
	return out
}
