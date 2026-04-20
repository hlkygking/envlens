package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envlens/internal/envdrift"
)

// EnvDriftTextReport renders a human-readable drift report.
func EnvDriftTextReport(results []envdrift.Result, summary envdrift.Summary) string {
	var sb strings.Builder

	sb.WriteString("=== Env Drift Report ===\n\n")

	drifted := filterDrift(results, envdrift.StatusDrifted)
	missing := filterDrift(results, envdrift.StatusMissing)

	if len(drifted) > 0 {
		sb.WriteString("[DRIFTED]\n")
		for _, r := range drifted {
			sb.WriteString(fmt.Sprintf("  %s (baseline: %q)\n", r.Key, r.Baseline))
			for env, val := range r.Values {
				if env == "baseline" {
					continue
				}
				sb.WriteString(fmt.Sprintf("    %s: %q\n", env, val))
			}
		}
		sb.WriteString("\n")
	}

	if len(missing) > 0 {
		sb.WriteString("[MISSING]\n")
		for _, r := range missing {
			sb.WriteString(fmt.Sprintf("  %s\n", r.Key))
			for env, val := range r.Values {
				if env == "baseline" {
					continue
				}
				if val == "" {
					sb.WriteString(fmt.Sprintf("    %s: <missing>\n", env))
				} else {
					sb.WriteString(fmt.Sprintf("    %s: %q\n", env, val))
				}
			}
		}
		sb.WriteString("\n")
	}

	if summary.Drifted == 0 && summary.Missing == 0 {
		sb.WriteString("No drift detected.\n\n")
	}

	sb.WriteString(fmt.Sprintf("Summary: total=%d match=%d drifted=%d missing=%d\n",
		summary.Total, summary.Match, summary.Drifted, summary.Missing))

	return sb.String()
}

// EnvDriftJSONReport renders drift results as a JSON string.
func EnvDriftJSONReport(results []envdrift.Result, summary envdrift.Summary) (string, error) {
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

func filterDrift(results []envdrift.Result, status envdrift.Status) []envdrift.Result {
	out := []envdrift.Result{}
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
