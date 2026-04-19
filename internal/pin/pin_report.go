package pin

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PinTextReport renders a human-readable report of pin check results.
func PinTextReport(entries []Entry) string {
	var sb strings.Builder

	matched := filterEntries(entries, StatusMatch)
	mismatched := filterEntries(entries, StatusMismatch)
	missing := filterEntries(entries, StatusMissing)

	if len(mismatched) > 0 {
		sb.WriteString("=== MISMATCHED ===\n")
		for _, e := range mismatched {
			fmt.Fprintf(&sb, "  %-30s expected=%q actual=%q\n", e.Key, e.Expected, e.Actual)
		}
		sb.WriteString("\n")
	}

	if len(missing) > 0 {
		sb.WriteString("=== MISSING ===\n")
		for _, e := range missing {
			fmt.Fprintf(&sb, "  %-30s expected=%q\n", e.Key, e.Expected)
		}
		sb.WriteString("\n")
	}

	if len(mismatched) == 0 && len(missing) == 0 {
		sb.WriteString("All pinned keys match.\n\n")
	}

	fmt.Fprintf(&sb, "Summary: %d matched, %d mismatched, %d missing\n",
		len(matched), len(mismatched), len(missing))

	return sb.String()
}

// PinJSONReport renders a JSON report of pin check results.
func PinJSONReport(entries []Entry) string {
	type output struct {
		Entries []Entry `json:"entries"`
		Summary struct {
			Matched    int `json:"matched"`
			Mismatched int `json:"mismatched"`
			Missing    int `json:"missing"`
		} `json:"summary"`
	}

	o := output{Entries: entries}
	o.Summary.Matched = len(filterEntries(entries, StatusMatch))
	o.Summary.Mismatched = len(filterEntries(entries, StatusMismatch))
	o.Summary.Missing = len(filterEntries(entries, StatusMissing))

	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

func filterEntries(entries []Entry, status string) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Status == status {
			out = append(out, e)
		}
	}
	return out
}
