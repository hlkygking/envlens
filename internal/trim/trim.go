package trim

import (
	"strings"
)

// Result holds the outcome of trimming a single key.
type Result struct {
	Key      string
	Original string
	Trimmed  string
	Changed  bool
}

// Summary holds aggregate trim statistics.
type Summary struct {
	Total   int
	Changed int
}

// Apply trims leading/trailing whitespace from all values in the map.
// If stripQuotes is true, surrounding single or double quotes are also removed.
func Apply(env map[string]string, stripQuotes bool) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		orig := v
		trimmed := strings.TrimSpace(v)
		if stripQuotes {
			trimmed = removeQuotes(trimmed)
		}
		results = append(results, Result{
			Key:      k,
			Original: orig,
			Trimmed:  trimmed,
			Changed:  orig != trimmed,
		})
	}
	return results
}

// ToMap converts trim results back to a plain map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Trimmed
	}
	return m
}

// GetSummary returns aggregate counts over the result set.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Changed {
			s.Changed++
		}
	}
	return s
}

// FilterChanged returns only the results where the value was modified.
func FilterChanged(results []Result) []Result {
	changed := make([]Result, 0)
	for _, r := range results {
		if r.Changed {
			changed = append(changed, r)
		}
	}
	return changed
}

func removeQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
