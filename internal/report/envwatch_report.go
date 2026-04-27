package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourorg/envlens/internal/envwatch"
)

// EnvWatchTextReport returns a human-readable watch diff report.
func EnvWatchTextReport(entries []envwatch.Entry, showUnchanged bool) string {
	var sb strings.Builder

	watched := filterWatchEntries(entries, showUnchanged)
	if len(watched) == 0 {
		sb.WriteString("No watched changes detected.\n")
		return sb.String()
	}

	sections := map[envwatch.Status][]envwatch.Entry{}
	for _, e := range watched {
		sections[e.Status] = append(sections[e.Status], e)
	}

	if s := sections[envwatch.StatusNew]; len(s) > 0 {
		sb.WriteString("=== New Keys ===\n")
		for _, e := range s {
			sb.WriteString(fmt.Sprintf("  + %s = %s\n", e.Key, e.NewValue))
		}
	}
	if s := sections[envwatch.StatusRemoved]; len(s) > 0 {
		sb.WriteString("=== Removed Keys ===\n")
		for _, e := range s {
			sb.WriteString(fmt.Sprintf("  - %s (was: %s)\n", e.Key, e.OldValue))
		}
	}
	if s := sections[envwatch.StatusChanged]; len(s) > 0 {
		sb.WriteString("=== Changed Keys ===\n")
		for _, e := range s {
			sb.WriteString(fmt.Sprintf("  ~ %s: %s -> %s\n", e.Key, e.OldValue, e.NewValue))
		}
	}
	if showUnchanged {
		if s := sections[envwatch.StatusUnchanged]; len(s) > 0 {
			sb.WriteString("=== Unchanged Keys ===\n")
			for _, e := range s {
				sb.WriteString(fmt.Sprintf("    %s = %s\n", e.Key, e.NewValue))
			}
		}
	}

	summary := envwatch.GetSummary(entries)
	sb.WriteString(fmt.Sprintf("\nSummary: %d new, %d removed, %d changed, %d unchanged\n",
		summary["new"], summary["removed"], summary["changed"], summary["unchanged"]))
	return sb.String()
}

// EnvWatchJSONReport returns a JSON-encoded watch report.
func EnvWatchJSONReport(entries []envwatch.Entry) (string, error) {
	type jsonEntry struct {
		Key      string `json:"key"`
		OldValue string `json:"old_value,omitempty"`
		NewValue string `json:"new_value,omitempty"`
		Status   string `json:"status"`
		Watched  bool   `json:"watched"`
	}
	out := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, jsonEntry{
			Key:      e.Key,
			OldValue: e.OldValue,
			NewValue: e.NewValue,
			Status:   string(e.Status),
			Watched:  e.Watched,
		})
	}
	summary := envwatch.GetSummary(entries)
	payload := map[string]any{"entries": out, "summary": summary}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterWatchEntries(entries []envwatch.Entry, showUnchanged bool) []envwatch.Entry {
	var out []envwatch.Entry
	for _, e := range entries {
		if !e.Watched {
			continue
		}
		if !showUnchanged && e.Status == envwatch.StatusUnchanged {
			continue
		}
		out = append(out, e)
	}
	return out
}
