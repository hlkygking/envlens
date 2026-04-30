package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/envlens/internal/envreplace"
)

// EnvReplaceTextReport renders a human-readable report of replace results.
func EnvReplaceTextReport(results []envreplace.Result, showUnchanged bool) string {
	var sb strings.Builder

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	replaced := filterReplace(results, envreplace.StatusReplaced)
	errored := filterReplace(results, envreplace.StatusError)

	if len(replaced) > 0 {
		sb.WriteString("=== Replaced ===\n")
		for _, r := range replaced {
			sb.WriteString(fmt.Sprintf("  %s: %q -> %q  (rule: %s)\n", r.Key, r.OldValue, r.NewValue, r.MatchedRule))
		}
	}

	if len(errored) > 0 {
		sb.WriteString("=== Errors ===\n")
		for _, r := range errored {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", r.Key, r.Error))
		}
	}

	if showUnchanged {
		unchanged := filterReplace(results, envreplace.StatusUnchanged)
		if len(unchanged) > 0 {
			sb.WriteString("=== Unchanged ===\n")
			for _, r := range unchanged {
				sb.WriteString(fmt.Sprintf("  %s=%q\n", r.Key, r.OldValue))
			}
		}
	}

	replacedN, unchangedN, erroredN := envreplace.GetSummary(results)
	sb.WriteString(fmt.Sprintf("\nSummary: %d replaced, %d unchanged, %d errors\n",
		replacedN, unchangedN, erroredN))

	return sb.String()
}

// EnvReplaceJSONReport renders a JSON report of replace results.
func EnvReplaceJSONReport(results []envreplace.Result) (string, error) {
	type entry struct {
		Key         string `json:"key"`
		OldValue    string `json:"old_value"`
		NewValue    string `json:"new_value"`
		Status      string `json:"status"`
		MatchedRule string `json:"matched_rule,omitempty"`
		Error       string `json:"error,omitempty"`
	}
	type summary struct {
		Replaced  int `json:"replaced"`
		Unchanged int `json:"unchanged"`
		Errored   int `json:"errored"`
	}
	type output struct {
		Entries []entry `json:"entries"`
		Summary summary `json:"summary"`
	}

	entries := make([]entry, 0, len(results))
	for _, r := range results {
		entries = append(entries, entry{
			Key:         r.Key,
			OldValue:    r.OldValue,
			NewValue:    r.NewValue,
			Status:      string(r.Status),
			MatchedRule: r.MatchedRule,
			Error:       r.Error,
		})
	}
	replacedN, unchangedN, erroredN := envreplace.GetSummary(results)
	out := output{
		Entries: entries,
		Summary: summary{Replaced: replacedN, Unchanged: unchangedN, Errored: erroredN},
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func filterReplace(results []envreplace.Result, status envreplace.Status) []envreplace.Result {
	var out []envreplace.Result
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
