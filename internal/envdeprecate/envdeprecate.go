package envdeprecate

import (
	"regexp"
	"strings"
)

// Status represents the deprecation state of a key.
type Status string

const (
	StatusDeprecated Status = "deprecated"
	StatusRenamed    Status = "renamed"
	StatusOK         Status = "ok"
)

// Rule defines a deprecation rule for a key pattern.
type Rule struct {
	Pattern     string
	Replacement string // empty means fully deprecated, no replacement
	Reason      string
}

// Result holds the outcome of checking a single key.
type Result struct {
	Key         string
	Value       string
	Status      Status
	Replacement string
	Reason      string
}

// Apply checks the given env map against the provided rules and returns results.
func Apply(env map[string]string, rules []Rule) []Result {
	var results []Result

	for key, val := range env {
		matched := false
		for _, rule := range rules {
			ok, err := regexp.MatchString("(?i)"+rule.Pattern, key)
			if err != nil || !ok {
				continue
			}
			matched = true
			status := StatusDeprecated
			if rule.Replacement != "" {
				status = StatusRenamed
			}
			results = append(results, Result{
				Key:         key,
				Value:       val,
				Status:      status,
				Replacement: rule.Replacement,
				Reason:      rule.Reason,
			})
			break
		}
		if !matched {
			results = append(results, Result{
				Key:    key,
				Value:  val,
				Status: StatusOK,
			})
		}
	}
	return results
}

// GetSummary returns counts of deprecated, renamed, and ok keys.
func GetSummary(results []Result) map[string]int {
	summary := map[string]int{
		"deprecated": 0,
		"renamed":    0,
		"ok":         0,
	}
	for _, r := range results {
		summary[strings.ToLower(string(r.Status))]++
	}
	return summary
}

// ParseRules parses rules from string slices in the format "PATTERN:REPLACEMENT:REASON".
// Replacement and Reason are optional.
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 3)
		rule := Rule{Pattern: strings.TrimSpace(parts[0])}
		if len(parts) >= 2 {
			rule.Replacement = strings.TrimSpace(parts[1])
		}
		if len(parts) >= 3 {
			rule.Reason = strings.TrimSpace(parts[2])
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
