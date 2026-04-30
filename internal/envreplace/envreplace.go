package envreplace

import (
	"fmt"
	"regexp"
	"strings"
)

// Strategy controls how replacement targets are matched.
type Strategy string

const (
	StrategyExact  Strategy = "exact"
	StrategyPrefix Strategy = "prefix"
	StrategyRegex  Strategy = "regex"
)

// Rule defines a single replacement operation.
type Rule struct {
	Target      string
	Replacement string
	Strategy    Strategy
}

// Status describes the outcome of applying a rule to a key.
type Status string

const (
	StatusReplaced  Status = "replaced"
	StatusUnchanged Status = "unchanged"
	StatusError     Status = "error"
)

// Result holds the outcome for a single environment variable.
type Result struct {
	Key        string
	OldValue   string
	NewValue   string
	Status     Status
	MatchedRule string
	Error      string
}

// Apply runs all rules against the provided env map and returns results.
func Apply(env map[string]string, rules []Rule) []Result {
	results := make([]Result, 0, len(env))
	for key, val := range env {
		res := Result{Key: key, OldValue: val, NewValue: val, Status: StatusUnchanged}
		for _, rule := range rules {
			newVal, matched, err := applyRule(val, rule)
			if err != nil {
				res.Status = StatusError
				res.Error = err.Error()
				res.MatchedRule = rule.Target
				break
			}
			if matched {
				res.NewValue = newVal
				res.Status = StatusReplaced
				res.MatchedRule = rule.Target
				break
			}
		}
		results = append(results, res)
	}
	return results
}

func applyRule(value string, rule Rule) (string, bool, error) {
	switch rule.Strategy {
	case StrategyExact:
		if value == rule.Target {
			return rule.Replacement, true, nil
		}
		return value, false, nil
	case StrategyPrefix:
		if strings.HasPrefix(value, rule.Target) {
			return rule.Replacement + strings.TrimPrefix(value, rule.Target), true, nil
		}
		return value, false, nil
	case StrategyRegex:
		re, err := regexp.Compile(rule.Target)
		if err != nil {
			return value, false, fmt.Errorf("invalid regex %q: %w", rule.Target, err)
		}
		newVal := re.ReplaceAllString(value, rule.Replacement)
		return newVal, newVal != value, nil
	default:
		if value == rule.Target {
			return rule.Replacement, true, nil
		}
		return value, false, nil
	}
}

// ToMap converts results back to a plain map, using NewValue for each key.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.NewValue
	}
	return m
}

// GetSummary returns counts of replaced, unchanged, and errored entries.
func GetSummary(results []Result) (replaced, unchanged, errored int) {
	for _, r := range results {
		switch r.Status {
		case StatusReplaced:
			replaced++
		case StatusUnchanged:
			unchanged++
		case StatusError:
			errored++
		}
	}
	return
}
