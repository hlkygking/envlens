package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envlens/internal/dedup"
)

// DedupTextReport returns a human-readable dedup report.
func DedupTextReport(results []dedup.Result) string {
	var sb strings.Builder
	total, dups := dedup.Summary(results)

	sb.WriteString("=== Deduplication Report ===\n")
	if dups == 0 {
		sb.WriteString("No duplicate keys found.\n")
	} else {
		sb.WriteString(fmt.Sprintf("Duplicates: %d / %d keys\n\n", dups, total))
		sb.WriteString("[DUPLICATES]\n")
		for _, r := range results {
			if !r.Duplicate {
				continue
			}
			sb.WriteString(fmt.Sprintf("  %-30s kept from %-15s sources: [%s]\n",
				r.Key, r.Kept, strings.Join(r.Sources, ", ")))
		}
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d total keys, %d duplicates resolved\n", total, dups))
	return sb.String()
}

// DedupJSONReport returns a JSON dedup report.
func DedupJSONReport(results []dedup.Result) (string, error) {
	total, dups := dedup.Summary(results)
	type row struct {
		Key       string   `json:"key"`
		Value     string   `json:"value"`
		Sources   []string `json:"sources"`
		Kept      string   `json:"kept"`
		Duplicate bool     `json:"duplicate"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{
			Key:       r.Key,
			Value:     r.Value,
			Sources:   r.Sources,
			Kept:      r.Kept,
			Duplicate: r.Duplicate,
		}
	}
	out := map[string]interface{}{
		"results":    rows,
		"total":      total,
		"duplicates": dups,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
