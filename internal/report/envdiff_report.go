package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/envdiff"
)

// EnvDiffTextReport renders a coloured text diff report.
func EnvDiffTextReport(entries []envdiff.Entry) string {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	var sb strings.Builder
	sb.WriteString("=== Env Diff ===\n")

	for _, e := range entries {
		if e.Change == envdiff.Unchanged {
			continue
		}
		sb.WriteString(envdiff.Format(e) + "\n")
	}

	s := envdiff.GetSummary(entries)
	sb.WriteString(fmt.Sprintf("\nSummary: +%d added, -%d removed, ~%d modified, %d unchanged\n",
		s.Added, s.Removed, s.Modified, s.Unchanged))
	return sb.String()
}

// EnvDiffJSONReport renders a JSON diff report.
func EnvDiffJSONReport(entries []envdiff.Entry) (string, error) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	type row struct {
		Key      string `json:"key"`
		Change   string `json:"change"`
		BaseVal  string `json:"base_value,omitempty"`
		TargetVal string `json:"target_value,omitempty"`
	}

	s := envdiff.GetSummary(entries)
	payload := struct {
		Entries []row `json:"entries"`
		Summary envdiff.Summary `json:"summary"`
	}{
		Summary: s,
	}

	for _, e := range entries {
		payload.Entries = append(payload.Entries, row{
			Key:       e.Key,
			Change:    string(e.Change),
			BaseVal:   e.BaseVal,
			TargetVal: e.TargetVal,
		})
	}

	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
