package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/envchain"
)

// EnvChainTextReport renders a human-readable report of chain results.
func EnvChainTextReport(results []envchain.Result, strategy envchain.Strategy) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== Envchain Report (strategy: %s) ===\n", strategy))

	if len(results) == 0 {
		sb.WriteString("No results.\n")
		return sb.String()
	}

	sorted := make([]envchain.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	overriddenCount := 0
	for _, r := range sorted {
		conflict := ""
		if r.Overridden {
			conflict = " [conflict]"
			overriddenCount++
		}
		sb.WriteString(fmt.Sprintf("  %-30s = %-20s (from: %s)%s\n", r.Key, r.Value, r.ResolvedBy, conflict))
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d keys resolved, %d conflicts\n", len(results), overriddenCount))
	return sb.String()
}

// EnvChainJSONReport renders a JSON report of chain results.
func EnvChainJSONReport(results []envchain.Result, strategy envchain.Strategy) (string, error) {
	type entry struct {
		Key        string `json:"key"`
		Value      string `json:"value"`
		ResolvedBy string `json:"resolved_by"`
		Overridden bool   `json:"overridden"`
	}
	type payload struct {
		Strategy string  `json:"strategy"`
		Total    int     `json:"total"`
		Conflicts int    `json:"conflicts"`
		Entries  []entry `json:"entries"`
	}

	entries := make([]entry, 0, len(results))
	conflicts := 0
	for _, r := range results {
		if r.Overridden {
			conflicts++
		}
		entries = append(entries, entry{Key: r.Key, Value: r.Value, ResolvedBy: r.ResolvedBy, Overridden: r.Overridden})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })

	p := payload{Strategy: string(strategy), Total: len(results), Conflicts: conflicts, Entries: entries}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
