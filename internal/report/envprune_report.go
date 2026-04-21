package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/envprune"
)

// EnvPruneTextReport renders a human-readable pruning report.
func EnvPruneTextReport(results []envprune.Result) string {
	var sb strings.Builder

	pruned := filterPrune(results, true)
	retained := filterPrune(results, false)

	sort.Slice(pruned, func(i, j int) bool { return pruned[i].Key < pruned[j].Key })
	sort.Slice(retained, func(i, j int) bool { return retained[i].Key < retained[j].Key })

	sb.WriteString("=== Pruned Keys ===\n")
	if len(pruned) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, r := range pruned {
			sb.WriteString(fmt.Sprintf("  - %s  [%s]\n", r.Key, r.Reason))
		}
	}

	sb.WriteString("\n=== Retained Keys ===\n")
	if len(retained) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, r := range retained {
			sb.WriteString(fmt.Sprintf("  + %s\n", r.Key))
		}
	}

	p, ret := envprune.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: %d pruned, %d retained\n", p, ret))
	return sb.String()
}

// EnvPruneJSONReport renders pruning results as JSON.
func EnvPruneJSONReport(results []envprune.Result) (string, error) {
	type row struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Pruned  bool   `json:"pruned"`
		Reason  string `json:"reason,omitempty"`
	}
	rows := make([]row, 0, len(results))
	for _, r := range results {
		rows = append(rows, row{Key: r.Key, Value: r.Value, Pruned: r.Pruned, Reason: r.Reason})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	p, ret := envprune.GetSummary(results)
	payload := map[string]interface{}{
		"results":  rows,
		"pruned":   p,
		"retained": ret,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterPrune(results []envprune.Result, pruned bool) []envprune.Result {
	var out []envprune.Result
	for _, r := range results {
		if r.Pruned == pruned {
			out = append(out, r)
		}
	}
	return out
}
