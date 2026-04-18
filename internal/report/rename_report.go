package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/envlens/internal/rename"
)

// RenameTextReport returns a human-readable report of rename results.
func RenameTextReport(results []rename.Result) string {
	var sb strings.Builder
	renamed, skipped, unchanged := rename.Summary(results)

	sb.WriteString("=== Rename Report ===\n")
	for _, r := range results {
		switch {
		case r.Skipped:
			fmt.Fprintf(&sb, "  [SKIPPED]   %s -> %s (key not found)\n", r.OldKey, r.NewKey)
		case r.Renamed:
			fmt.Fprintf(&sb, "  [RENAMED]   %s -> %s\n", r.OldKey, r.NewKey)
		default:
			fmt.Fprintf(&sb, "  [UNCHANGED] %s\n", r.OldKey)
		}
	}
	sb.WriteString("\n--- Summary ---\n")
	fmt.Fprintf(&sb, "  Renamed:   %d\n", renamed)
	fmt.Fprintf(&sb, "  Skipped:   %d\n", skipped)
	fmt.Fprintf(&sb, "  Unchanged: %d\n", unchanged)
	return sb.String()
}

type renameJSONEntry struct {
	OldKey  string `json:"old_key"`
	NewKey  string `json:"new_key"`
	Status  string `json:"status"`
}

type renameJSONReport struct {
	Results  []renameJSONEntry `json:"results"`
	Renamed  int               `json:"renamed"`
	Skipped  int               `json:"skipped"`
	Unchanged int              `json:"unchanged"`
}

// RenameJSONReport returns a JSON-encoded report of rename results.
func RenameJSONReport(results []rename.Result) (string, error) {
	var entries []renameJSONEntry
	for _, r := range results {
		status := "unchanged"
		if r.Renamed {
			status = "renamed"
		} else if r.Skipped {
			status = "skipped"
		}
		entries = append(entries, renameJSONEntry{OldKey: r.OldKey, NewKey: r.NewKey, Status: status})
	}
	renamed, skipped, unchanged := rename.Summary(results)
	rep := renameJSONReport{Results: entries, Renamed: renamed, Skipped: skipped, Unchanged: unchanged}
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
