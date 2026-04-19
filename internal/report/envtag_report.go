package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envlens/internal/envtag"
)

// EnvTagTextReport renders a human-readable tag report.
func EnvTagTextReport(results []envtag.Result) string {
	var sb strings.Builder
	tagged, untagged := envtag.Summary(results)

	sb.WriteString("=== Env Tag Report ===\n")
	sb.WriteString(fmt.Sprintf("Tagged: %d | Untagged: %d\n\n", tagged, untagged))

	sb.WriteString("[Tagged]\n")
	hasTagged := false
	for _, r := range results {
		if r.Tagged {
			hasTagged = true
			sb.WriteString(fmt.Sprintf("  %-30s %s\n", r.Key, strings.Join(r.Tags, ", ")))
		}
	}
	if !hasTagged {
		sb.WriteString("  (none)\n")
	}

	sb.WriteString("\n[Untagged]\n")
	hasUntagged := false
	for _, r := range results {
		if !r.Tagged {
			hasUntagged = true
			sb.WriteString(fmt.Sprintf("  %s\n", r.Key))
		}
	}
	if !hasUntagged {
		sb.WriteString("  (none)\n")
	}
	return sb.String()
}

// EnvTagJSONReport renders a JSON tag report.
func EnvTagJSONReport(results []envtag.Result) (string, error) {
	tagged, untagged := envtag.Summary(results)
	payload := map[string]interface{}{
		"summary": map[string]int{
			"tagged":   tagged,
			"untagged": untagged,
		},
		"results": results,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
