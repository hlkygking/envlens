package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/envnorm"
)

// EnvNormTextReport returns a human-readable normalization report.
func EnvNormTextReport(results []envnorm.Result) string {
	var sb strings.Builder
	summary := envnorm.GetSummary(results)

	sorted := make([]envnorm.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OriginalKey < sorted[j].OriginalKey
	})

	sb.WriteString("=== Env Normalization ===\n")
	changed := filterNorm(sorted, true)
	if len(changed) > 0 {
		sb.WriteString("\n[CHANGED]\n")
		for _, r := range changed {
			fmt.Fprintf(&sb, "  %s -> %s\n", r.OriginalKey, r.NormalizedKey)
		}
	}

	unchanged := filterNorm(sorted, false)
	if len(unchanged) > 0 {
		sb.WriteString("\n[UNCHANGED]\n")
		for _, r := range unchanged {
			fmt.Fprintf(&sb, "  %s\n", r.OriginalKey)
		}
	}

	fmt.Fprintf(&sb, "\nSummary: %d total, %d changed, %d unchanged\n",
		summary.Total, summary.Changed, summary.Unchanged)
	return sb.String()
}

// EnvNormJSONReport returns a JSON normalization report.
func EnvNormJSONReport(results []envnorm.Result) (string, error) {
	summary := envnorm.GetSummary(results)
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

func filterNorm(results []envnorm.Result, changed bool) []envnorm.Result {
	var out []envnorm.Result
	for _, r := range results {
		if r.Changed == changed {
			out = append(out, r)
		}
	}
	return out
}
