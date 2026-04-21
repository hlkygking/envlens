package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/envcompact"
)

// EnvCompactTextReport renders a human-readable compaction report.
func EnvCompactTextReport(results []envcompact.Result) string {
	var sb strings.Builder

	kept, removed := envcompact.GetSummary(results)
	fmt.Fprintf(&sb, "Compaction Report\n")
	fmt.Fprintf(&sb, "  Kept:    %d\n", kept)
	fmt.Fprintf(&sb, "  Removed: %d\n\n", removed)

	removed_entries := filterCompact(results, true)
	if len(removed_entries) == 0 {
		sb.WriteString("No entries removed.\n")
		return sb.String()
	}

	sb.WriteString("Removed:\n")
	sort.Slice(removed_entries, func(i, j int) bool {
		return removed_entries[i].Key < removed_entries[j].Key
	})
	for _, r := range removed_entries {
		fmt.Fprintf(&sb, "  - %-30s (%s)\n", r.Key, r.Reason)
	}
	return sb.String()
}

// EnvCompactJSONReport renders a JSON compaction report.
func EnvCompactJSONReport(results []envcompact.Result) (string, error) {
	kept, removed := envcompact.GetSummary(results)

	type row struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Removed bool   `json:"removed"`
		Reason  string `json:"reason,omitempty"`
	}

	rows := make([]row, 0, len(results))
	for _, r := range results {
		rows = append(rows, row{Key: r.Key, Value: r.Value, Removed: r.Removed, Reason: r.Reason})
	}

	out := map[string]interface{}{
		"summary": map[string]int{"kept": kept, "removed": removed},
		"entries": rows,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterCompact(results []envcompact.Result, removed bool) []envcompact.Result {
	out := make([]envcompact.Result, 0)
	for _, r := range results {
		if r.Removed == removed {
			out = append(out, r)
		}
	}
	return out
}
