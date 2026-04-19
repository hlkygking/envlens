package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/envlens/internal/envgroup"
)

// EnvGroupTextReport renders a human-readable group report.
func EnvGroupTextReport(r envgroup.Result) string {
	var sb strings.Builder
	sb.WriteString("=== Env Groups ===\n")
	for _, g := range r.Groups {
		sb.WriteString(fmt.Sprintf("\n[%s] (pattern: %s) — %d key(s)\n", g.Name, g.Pattern, len(g.Keys)))
		for _, k := range g.Keys {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	if len(r.Ungrouped) > 0 {
		sb.WriteString(fmt.Sprintf("\n[ungrouped] — %d key(s)\n", len(r.Ungrouped)))
		for _, k := range r.Ungrouped {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	sm := envgroup.Summary(r)
	sb.WriteString("\n--- Summary ---\n")
	for k, v := range sm {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", k, v))
	}
	return sb.String()
}

// EnvGroupJSONReport renders a JSON group report.
func EnvGroupJSONReport(r envgroup.Result) (string, error) {
	type payload struct {
		Groups    []envgroup.Group `json:"groups"`
		Ungrouped []string         `json:"ungrouped"`
		Summary   map[string]int   `json:"summary"`
	}
	p := payload{
		Groups:    r.Groups,
		Ungrouped: r.Ungrouped,
		Summary:   envgroup.Summary(r),
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
