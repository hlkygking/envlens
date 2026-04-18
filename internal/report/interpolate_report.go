package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourorg/envlens/internal/interpolate"
)

// InterpolateTextReport renders interpolation results as plain text.
func InterpolateTextReport(results []interpolate.Result, cycles []string) string {
	var sb strings.Builder
	summary := interpolate.GetSummary(results)

	sb.WriteString("=== Interpolation Report ===\n")
	sb.WriteString(fmt.Sprintf("Total: %d | Resolved: %d | Missing refs: %d\n\n",
		summary.Total, summary.Resolved, summary.Missing))

	if len(cycles) > 0 {
		sb.WriteString("[CYCLES DETECTED]\n")
		for _, c := range cycles {
			sb.WriteString(fmt.Sprintf("  !! %s\n", c))
		}
		sb.WriteString("\n")
	}

	for _, r := range results {
		if len(r.Refs) == 0 {
			continue
		}
		status := "OK"
		if len(r.Missing) > 0 {
			status = fmt.Sprintf("MISSING(%s)", strings.Join(r.Missing, ","))
		}
		sb.WriteString(fmt.Sprintf("  %s: %s -> %s [%s]\n", r.Key, r.Original, r.Resolved, status))
	}

	if summary.Resolved == 0 && summary.Missing == 0 {
		sb.WriteString("  No interpolation references found.\n")
	}
	return sb.String()
}

// InterpolateJSONReport renders interpolation results as JSON.
func InterpolateJSONReport(results []interpolate.Result, cycles []string) (string, error) {
	type jsonResult struct {
		Key      string   `json:"key"`
		Original string   `json:"original"`
		Resolved string   `json:"resolved"`
		Refs     []string `json:"refs"`
		Missing  []string `json:"missing"`
	}
	summary := interpolate.GetSummary(results)
	payload := struct {
		Summary interpolate.Summary `json:"summary"`
		Cycles  []string            `json:"cycles"`
		Results []jsonResult        `json:"results"`
	}{
		Summary: summary,
		Cycles:  cycles,
	}
	for _, r := range results {
		payload.Results = append(payload.Results, jsonResult{
			Key:      r.Key,
			Original: r.Original,
			Resolved: r.Resolved,
			Refs:     r.Refs,
			Missing:  r.Missing,
		})
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
