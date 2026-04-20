package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/envcount"
)

// EnvCountTextReport renders a human-readable summary of env value length analysis.
func EnvCountTextReport(results []envcount.Result) string {
	var sb strings.Builder

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	sb.WriteString("=== Env Count Report ===\n")
	for _, r := range results {
		sensTag := ""
		if r.Sensitive {
			sensTag = " [sensitive]"
		}
		sb.WriteString(fmt.Sprintf("  %-30s len=%-5d status=%-8s%s\n",
			r.Key, r.Length, r.Status, sensTag))
	}

	summary := envcount.GetSummary(results)
	sb.WriteString("\n--- Summary ---\n")
	sb.WriteString(fmt.Sprintf("  Total:     %d\n", summary.Total))
	sb.WriteString(fmt.Sprintf("  Empty:     %d\n", summary.Empty))
	sb.WriteString(fmt.Sprintf("  Short:     %d\n", summary.Short))
	sb.WriteString(fmt.Sprintf("  Medium:    %d\n", summary.Medium))
	sb.WriteString(fmt.Sprintf("  Long:      %d\n", summary.Long))
	sb.WriteString(fmt.Sprintf("  Sensitive: %d\n", summary.Sensitive))

	return sb.String()
}

type envCountJSON struct {
	Key       string `json:"key"`
	Length    int    `json:"length"`
	Status    string `json:"status"`
	Sensitive bool   `json:"sensitive"`
}

type envCountReport struct {
	Entries []envCountJSON   `json:"entries"`
	Summary envcount.Summary `json:"summary"`
}

// EnvCountJSONReport renders a JSON report of env value length analysis.
func EnvCountJSONReport(results []envcount.Result) (string, error) {
	entries := make([]envCountJSON, 0, len(results))
	for _, r := range results {
		entries = append(entries, envCountJSON{
			Key:       r.Key,
			Length:    r.Length,
			Status:    string(r.Status),
			Sensitive: r.Sensitive,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	report := envCountReport{
		Entries: entries,
		Summary: envcount.GetSummary(results),
	}
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
