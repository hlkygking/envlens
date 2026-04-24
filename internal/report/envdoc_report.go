package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/envdoc"
)

// EnvDocTextReport renders documented/undocumented env vars as plain text.
func EnvDocTextReport(results []envdoc.Result, showAll bool) string {
	var sb strings.Builder

	sorted := make([]envdoc.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	sb.WriteString("=== Documented ===\n")
	for _, r := range sorted {
		if !r.Found {
			continue
		}
		line := fmt.Sprintf("  %-30s %s", r.Key, r.Description)
		if r.Default != "" {
			line += fmt.Sprintf(" (default: %s)", r.Default)
		}
		if r.Required {
			line += " [required]"
		}
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n=== Undocumented ===\n")
	for _, r := range sorted {
		if r.Found {
			continue
		}
		sb.WriteString(fmt.Sprintf("  %s\n", r.Key))
	}

	s := envdoc.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: total=%d documented=%d undocumented=%d\n",
		s.Total, s.Documented, s.Undocumented))
	return sb.String()
}

// EnvDocJSONReport renders results as JSON.
func EnvDocJSONReport(results []envdoc.Result) (string, error) {
	type jsonEntry struct {
		Key         string `json:"key"`
		Documented  bool   `json:"documented"`
		Description string `json:"description,omitempty"`
		Example     string `json:"example,omitempty"`
		Required    bool   `json:"required"`
		Default     string `json:"default,omitempty"`
	}

	entries := make([]jsonEntry, 0, len(results))
	for _, r := range results {
		entries = append(entries, jsonEntry{
			Key:         r.Key,
			Documented:  r.Found,
			Description: r.Description,
			Example:     r.Example,
			Required:    r.Required,
			Default:     r.Default,
		})
	}

	s := envdoc.GetSummary(results)
	out := map[string]any{
		"entries": entries,
		"summary": s,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
