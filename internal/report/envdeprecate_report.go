package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envlens/internal/envdeprecate"
)

// EnvDeprecateTextReport renders a human-readable deprecation report.
func EnvDeprecateTextReport(results []envdeprecate.Result, showOK bool) string {
	var sb strings.Builder

	deprecated := filterDeprecate(results, envdeprecate.StatusDeprecated)
	renamed := filterDeprecate(results, envdeprecate.StatusRenamed)
	ok := filterDeprecate(results, envdeprecate.StatusOK)

	if len(deprecated) > 0 {
		sb.WriteString("=== Deprecated ===\n")
		for _, r := range deprecated {
			line := fmt.Sprintf("  [DEPRECATED] %s", r.Key)
			if r.Reason != "" {
				line += fmt.Sprintf(" — %s", r.Reason)
			}
			sb.WriteString(line + "\n")
		}
	}

	if len(renamed) > 0 {
		sb.WriteString("=== Renamed ===\n")
		for _, r := range renamed {
			line := fmt.Sprintf("  [RENAMED] %s → %s", r.Key, r.Replacement)
			if r.Reason != "" {
				line += fmt.Sprintf(" (%s)", r.Reason)
			}
			sb.WriteString(line + "\n")
		}
	}

	if showOK && len(ok) > 0 {
		sb.WriteString("=== OK ===\n")
		for _, r := range ok {
			sb.WriteString(fmt.Sprintf("  [OK] %s\n", r.Key))
		}
	}

	summary := envdeprecate.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: %d deprecated, %d renamed, %d ok\n",
		summary["deprecated"], summary["renamed"], summary["ok"]))

	return sb.String()
}

// EnvDeprecateJSONReport renders a JSON deprecation report.
func EnvDeprecateJSONReport(results []envdeprecate.Result) (string, error) {
	summary := envdeprecate.GetSummary(results)
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

func filterDeprecate(results []envdeprecate.Result, status envdeprecate.Status) []envdeprecate.Result {
	var out []envdeprecate.Result
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
