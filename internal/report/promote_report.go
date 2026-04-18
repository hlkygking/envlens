package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/your-org/envlens/internal/promote"
)

// PromoteTextReport renders promotion results as human-readable text.
func PromoteTextReport(results []promote.Result, summary promote.Summary) string {
	var sb strings.Builder

	promoted := filterPromote(results, "promoted")
	skipped := filterPromote(results, "skipped")

	if len(promoted) > 0 {
		sb.WriteString("PROMOTED:\n")
		for _, r := range promoted {
			line := fmt.Sprintf("  + %s = %s", r.Key, r.Value)
			if r.Message != "" {
				line += fmt.Sprintf(" (%s)", r.Message)
			}
			sb.WriteString(line + "\n")
		}
	}

	if len(skipped) > 0 {
		sb.WriteString("SKIPPED:\n")
		for _, r := range skipped {
			line := fmt.Sprintf("  ~ %s", r.Key)
			if r.Message != "" {
				line += fmt.Sprintf(" (%s)", r.Message)
			}
			sb.WriteString(line + "\n")
		}
	}

	if len(promoted) == 0 && len(skipped) == 0 {
		sb.WriteString("No changes.\n")
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d promoted, %d skipped, %d conflicts\n",
		summary.Promoted, summary.Skipped, summary.Conflicts))

	return sb.String()
}

// PromoteJSONReport renders promotion results as JSON.
func PromoteJSONReport(results []promote.Result, summary promote.Summary) (string, error) {
	payload := map[string]interface{}{
		"results": results,
		"summary": summary,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterPromote(results []promote.Result, status string) []promote.Result {
	var out []promote.Result
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
