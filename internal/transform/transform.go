package transform

import (
	"fmt"
	"strings"
)

// Op represents a transformation operation type.
type Op string

const (
	OpUppercase Op = "uppercase"
	OpLowercase Op = "lowercase"
	OpTrimSpace Op = "trimspace"
	OpPrefix    Op = "prefix"
	OpSuffix    Op = "suffix"
	OpReplace   Op = "replace"
)

// Rule defines a single transformation to apply.
type Rule struct {
	Op    Op
	Arg   string // used by prefix, suffix, replace
	Arg2  string // used by replace (replacement value)
}

// Result holds the outcome of transforming a single key.
type Result struct {
	Key      string
	Original string
	Value    string
	Applied  []Op
}

// Apply applies a list of rules to each entry in the env map.
func Apply(env map[string]string, rules []Rule) ([]Result, error) {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		r := Result{Key: k, Original: v, Value: v}
		for _, rule := range rules {
			transformed, err := applyRule(r.Value, rule)
			if err != nil {
				return nil, fmt.Errorf("key %s: %w", k, err)
			}
			r.Value = transformed
			r.Applied = append(r.Applied, rule.Op)
		}
		results = append(results, r)
	}
	return results, nil
}

func applyRule(val string, rule Rule) (string, error) {
	switch rule.Op {
	case OpUppercase:
		return strings.ToUpper(val), nil
	case OpLowercase:
		return strings.ToLower(val), nil
	case OpTrimSpace:
		return strings.TrimSpace(val), nil
	case OpPrefix:
		return rule.Arg + val, nil
	case OpSuffix:
		return val + rule.Arg, nil
	case OpReplace:
		return strings.ReplaceAll(val, rule.Arg, rule.Arg2), nil
	default:
		return "", fmt.Errorf("unknown op: %s", rule.Op)
	}
}

// ToMap converts results back to a plain map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}

// Summary returns counts of transformed vs unchanged entries.
func Summary(results []Result) (transformed, unchanged int) {
	for _, r := range results {
		if r.Value != r.Original {
			transformed++
		} else {
			unchanged++
		}
	}
	return
}
