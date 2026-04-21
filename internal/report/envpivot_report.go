package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/envpivot"
)

// EnvPivotTextReport renders a human-readable pivot table.
func EnvPivotTextReport(entries []envpivot.Entry, scopes []string) string {
	var sb strings.Builder
	summary := envpivot.GetSummary(entries)

	sb.WriteString("=== Env Pivot Report ===\n")
	if len(entries) == 0 {
		sb.WriteString("No entries.\n")
		return sb.String()
	}

	// Header row
	sortedScopes := make([]string, len(scopes))
	copy(sortedScopes, scopes)
	sort.Strings(sortedScopes)

	sb.WriteString(fmt.Sprintf("%-30s", "KEY"))
	for _, sc := range sortedScopes {
		sb.WriteString(fmt.Sprintf(" %-20s", sc))
	}
	sb.WriteString(" STATUS\n")
	sb.WriteString(strings.Repeat("-", 30+21*len(sortedScopes)+8) + "\n")

	for _, e := range entries {
		status := "uniform"
		if !e.Uniform {
			status = "DIVERGENT"
		}
		sb.WriteString(fmt.Sprintf("%-30s", e.Key))
		for _, sc := range sortedScopes {
			val, ok := e.Values[sc]
			if !ok {
				val = "(missing)"
			}
			if len(val) > 18 {
				val = val[:15] + "..."
			}
			sb.WriteString(fmt.Sprintf(" %-20s", val))
		}
		sb.WriteString(fmt.Sprintf(" %s\n", status))
	}

	sb.WriteString(fmt.Sprintf("\nSummary: total=%d uniform=%d divergent=%d missing_in_any=%d\n",
		summary.TotalKeys, summary.UniformKeys, summary.DivergentKeys, summary.MissingInAny))
	return sb.String()
}

// EnvPivotJSONReport renders entries and summary as JSON.
func EnvPivotJSONReport(entries []envpivot.Entry) (string, error) {
	summary := envpivot.GetSummary(entries)
	payload := map[string]any{
		"entries": entries,
		"summary": summary,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
