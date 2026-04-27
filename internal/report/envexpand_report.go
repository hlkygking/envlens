package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/envexpand"
)

// EnvExpandTextReport returns a human-readable report of expansion results.
func EnvExpandTextReport(results []envexpand.Result, showUnchanged bool) string {
	var sb strings.Builder

	sorted := make([]envexpand.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	expanded := filterExpand(sorted, "ok")
	unresolved := filterExpand(sorted, "unresolved")
	unchanged := filterExpand(sorted, "unchanged")

	if len(expanded) > 0 {
		sb.WriteString("=== Expanded ===\n")
		for _, r := range expanded {
			sb.WriteString(fmt.Sprintf("  %s: %s -> %s\n", r.Key, r.Original, r.Expanded))
		}
	}

	if len(unresolved) > 0 {
		sb.WriteString("=== Unresolved ===\n")
		for _, r := range unresolved {
			sb.WriteString(fmt.Sprintf("  %s: %s (refs: %s)\n", r.Key, r.Original, strings.Join(r.Refs, ", ")))
		}
	}

	if showUnchanged && len(unchanged) > 0 {
		sb.WriteString("=== Unchanged ===\n")
		for _, r := range unchanged {
			sb.WriteString(fmt.Sprintf("  %s=%s\n", r.Key, r.Value()))
		}
	}

	s := envexpand.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: total=%d expanded=%d unresolved=%d unchanged=%d\n",
		s.Total, s.Expanded, s.Unresolved, s.Unchanged))

	return sb.String()
}

// EnvExpandJSONReport returns a JSON-encoded report of expansion results.
func EnvExpandJSONReport(results []envexpand.Result) (string, error) {
	type entry struct {
		Key      string   `json:"key"`
		Original string   `json:"original"`
		Expanded string   `json:"expanded"`
		Refs     []string `json:"refs,omitempty"`
		Status   string   `json:"status"`
	}

	entries := make([]entry, len(results))
	for i, r := range results {
		entries[i] = entry{
			Key:      r.Key,
			Original: r.Original,
			Expanded: r.Expanded,
			Refs:     r.Refs,
			Status:   r.Status,
		}
	}

	s := envexpand.GetSummary(results)
	out := map[string]interface{}{
		"results": entries,
		"summary": s,
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterExpand(results []envexpand.Result, status string) []envexpand.Result {
	var out []envexpand.Result
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
