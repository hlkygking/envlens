package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/scope"
)

// ScopeTextReport renders a human-readable cross-scope comparison.
func ScopeTextReport(results []scope.Result) string {
	var sb strings.Builder
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })

	uniform, divergent := scope.Summary(results)
	sb.WriteString(fmt.Sprintf("Scope Comparison: %d keys (%d uniform, %d divergent)\n", len(results), uniform, divergent))
	sb.WriteString(strings.Repeat("-", 48) + "\n")

	for _, r := range results {
		status := "UNIFORM"
		if !r.Uniform {
			status = "DIVERGENT"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", status, r.Key))
		scopes := make([]string, 0, len(r.Values))
		for s := range r.Values {
			scopes = append(scopes, s)
		}
		sort.Strings(scopes)
		for _, s := range scopes {
			sb.WriteString(fmt.Sprintf("  %s = %s\n", s, r.Values[s]))
		}
	}
	return sb.String()
}

// ScopeJSONReport renders results as JSON.
func ScopeJSONReport(results []scope.Result) (string, error) {
	type entry struct {
		Key     string            `json:"key"`
		Uniform bool              `json:"uniform"`
		Values  map[string]string `json:"values"`
	}
	uniform, divergent := scope.Summary(results)
	payload := struct {
		Uniform  int     `json:"uniform"`
		Divergent int    `json:"divergent"`
		Keys     []entry `json:"keys"`
	}{
		Uniform:  uniform,
		Divergent: divergent,
	}
	for _, r := range results {
		payload.Keys = append(payload.Keys, entry{Key: r.Key, Uniform: r.Uniform, Values: r.Values})
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
