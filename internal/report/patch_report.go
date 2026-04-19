package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourorg/envlens/internal/patch"
)

// PatchTextReport renders patch results as human-readable text.
func PatchTextReport(results []patch.Result) string {
	var sb strings.Builder
	applied, skipped := patch.Summary(results)
	sb.WriteString("=== Patch Results ===\n")

	for _, r := range results {
		status := "OK"
		if !r.Applied {
			status = "SKIP"
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s %s — %s\n", status, strings.ToUpper(string(r.Rule.Op)), r.Rule.Key, r.Note))
	}

	if len(results) == 0 {
		sb.WriteString("  (no rules applied)\n")
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d applied, %d skipped\n", applied, skipped))
	return sb.String()
}

// PatchJSONReport renders patch results as JSON.
func PatchJSONReport(results []patch.Result) (string, error) {
	type entry struct {
		Op      string `json:"op"`
		Key     string `json:"key"`
		To      string `json:"to,omitempty"`
		Value   string `json:"value,omitempty"`
		Applied bool   `json:"applied"`
		Note    string `json:"note"`
	}
	applied, skipped := patch.Summary(results)
	type payload struct {
		Results []entry `json:"results"`
		Applied int     `json:"applied"`
		Skipped int     `json:"skipped"`
	}
	p := payload{Applied: applied, Skipped: skipped}
	for _, r := range results {
		p.Results = append(p.Results, entry{
			Op:      string(r.Rule.Op),
			Key:     r.Rule.Key,
			To:      r.Rule.To,
			Value:   r.Rule.Value,
			Applied: r.Applied,
			Note:    r.Note,
		})
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
