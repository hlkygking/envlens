package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/envlens/internal/sort"
)

// SortTextReport renders sorted entries as human-readable text.
func SortTextReport(res sort.Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Sorted %d entries (order: %s)\n", res.Total, res.Order)
	fmt.Fprintln(&sb, strings.Repeat("-", 40))

	currentGroup := ""
	for _, e := range res.Entries {
		if string(res.Order) == "group" && e.Group != currentGroup {
			if currentGroup != "" {
				fmt.Fprintln(&sb)
			}
			if e.Group != "" {
				fmt.Fprintf(&sb, "[%s]\n", e.Group)
			}
			currentGroup = e.Group
		}
		fmt.Fprintf(&sb, "  %s=%s\n", e.Key, e.Value)
	}

	return sb.String()
}

// SortJSONReport renders sorted entries as JSON.
func SortJSONReport(res sort.Result) string {
	type jsonEntry struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Group string `json:"group,omitempty"`
	}
	type jsonReport struct {
		Order   string      `json:"order"`
		Total   int         `json:"total"`
		Entries []jsonEntry `json:"entries"`
	}

	entries := make([]jsonEntry, len(res.Entries))
	for i, e := range res.Entries {
		entries[i] = jsonEntry{Key: e.Key, Value: e.Value, Group: e.Group}
	}

	out := jsonReport{
		Order:   string(res.Order),
		Total:   res.Total,
		Entries: entries,
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return `{"error":"marshal failed"}`
	}
	return string(b)
}
