package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/envrotate"
)

// EnvRotateTextReport renders a human-readable rotation report.
func EnvRotateTextReport(results []envrotate.Result, dryRun bool) string {
	var sb strings.Builder

	mode := ""
	if dryRun {
		mode = " [DRY RUN]"
	}
	sb.WriteString(fmt.Sprintf("=== Env Rotate Report%s ===\n", mode))

	rotated := filterRotate(results, true)
	unchanged := filterRotate(results, false)

	if len(rotated) > 0 {
		sb.WriteString("\n[ROTATED]\n")
		sort.Slice(rotated, func(i, j int) bool { return rotated[i].Key < rotated[j].Key })
		for _, r := range rotated {
			sb.WriteString(fmt.Sprintf("  %-30s %s -> %s\n", r.Key, r.OldValue, r.NewValue))
		}
	}

	if len(unchanged) > 0 {
		sb.WriteString("\n[UNCHANGED]\n")
		sort.Slice(unchanged, func(i, j int) bool { return unchanged[i].Key < unchanged[j].Key })
		for _, r := range unchanged {
			sb.WriteString(fmt.Sprintf("  %-30s (skipped)\n", r.Key))
		}
	}

	summary := envrotate.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: %d total, %d rotated, %d unchanged\n",
		summary["total"], summary["rotated"], summary["unchanged"]))

	return sb.String()
}

// EnvRotateJSONReport renders a JSON rotation report.
func EnvRotateJSONReport(results []envrotate.Result, dryRun bool) (string, error) {
	type payload struct {
		DryRun  bool                 `json:"dry_run"`
		Results []envrotate.Result   `json:"results"`
		Summary map[string]int       `json:"summary"`
	}
	p := payload{
		DryRun:  dryRun,
		Results: results,
		Summary: envrotate.GetSummary(results),
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterRotate(results []envrotate.Result, rotated bool) []envrotate.Result {
	out := make([]envrotate.Result, 0)
	for _, r := range results {
		if r.Rotated == rotated {
			out = append(out, r)
		}
	}
	return out
}
