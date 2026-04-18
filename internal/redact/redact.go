package redact

import (
	"strings"
)

// Rule defines a redaction rule for a specific key pattern.
type Rule struct {
	KeyPattern string
	Replacement string
}

// Result holds the outcome of redacting a single key-value pair.
type Result struct {
	Key       string
	Original  string
	Redacted  string
	WasRedacted bool
}

var defaultRules = []Rule{
	{KeyPattern: "PASSWORD", Replacement: "[REDACTED]"},
	{KeyPattern: "SECRET", Replacement: "[REDACTED]"},
	{KeyPattern: "TOKEN", Replacement: "[REDACTED]"},
	{KeyPattern: "API_KEY", Replacement: "[REDACTED]"},
	{KeyPattern: "PRIVATE", Replacement: "[REDACTED]"},
	{KeyPattern: "CREDENTIAL", Replacement: "[REDACTED]"},
}

// Apply redacts values in the provided map using the given rules.
// If rules is nil, default rules are used.
func Apply(env map[string]string, rules []Rule) []Result {
	if rules == nil {
		rules = defaultRules
	}

	results := make([]Result, 0, len(env))
	for k, v := range env {
		r := Result{Key: k, Original: v, Redacted: v, WasRedacted: false}
		upper := strings.ToUpper(k)
		for _, rule := range rules {
			if strings.Contains(upper, strings.ToUpper(rule.KeyPattern)) {
				r.Redacted = rule.Replacement
				r.WasRedacted = true
				break
			}
		}
		results = append(results, r)
	}
	return results
}

// ToMap converts a slice of Results into a redacted map.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		out[r.Key] = r.Redacted
	}
	return out
}

// Summary returns counts of redacted vs total entries.
func Summary(results []Result) (total int, redacted int) {
	total = len(results)
	for _, r := range results {
		if r.WasRedacted {
			redacted++
		}
	}
	return
}
