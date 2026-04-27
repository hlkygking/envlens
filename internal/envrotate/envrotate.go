package envrotate

import (
	"fmt"
	"strings"
)

// Strategy defines how rotation generates new values.
type Strategy string

const (
	StrategyRedact  Strategy = "redact"
	StrategyIncrement Strategy = "increment"
	StrategyBlank   Strategy = "blank"
)

// Result holds the outcome of rotating a single key.
type Result struct {
	Key      string
	OldValue string
	NewValue string
	Rotated  bool
	Reason   string
}

// Options controls rotation behaviour.
type Options struct {
	Strategy Strategy
	Keys     []string // explicit keys to rotate; empty = rotate all sensitive
	DryRun   bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy: StrategyRedact,
	}
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, kw := range []string{"secret", "password", "passwd", "token", "key", "api", "auth", "credential", "private"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func shouldRotate(key string, explicit []string) bool {
	if len(explicit) == 0 {
		return isSensitive(key)
	}
	for _, k := range explicit {
		if k == key {
			return true
		}
	}
	return false
}

func rotateValue(key, old string, s Strategy) string {
	switch s {
	case StrategyIncrement:
		return fmt.Sprintf("%s_rotated", old)
	case StrategyBlank:
		return ""
	default: // redact
		return fmt.Sprintf("REDACTED_%s", strings.ToUpper(key))
	}
}

// Apply rotates values in env according to opts.
func Apply(env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(env))
	for key, val := range env {
		if !shouldRotate(key, opts.Keys) {
			results = append(results, Result{Key: key, OldValue: val, NewValue: val, Rotated: false, Reason: "not targeted"})
			continue
		}
		newVal := rotateValue(key, val, opts.Strategy)
		results = append(results, Result{Key: key, OldValue: val, NewValue: newVal, Rotated: true, Reason: string(opts.Strategy)})
	}
	return results
}

// ToMap converts results to a plain map (new values).
func ToMap(results []Result) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		out[r.Key] = r.NewValue
	}
	return out
}

// GetSummary returns counts of rotated vs unchanged.
func GetSummary(results []Result) map[string]int {
	summary := map[string]int{"total": len(results), "rotated": 0, "unchanged": 0}
	for _, r := range results {
		if r.Rotated {
			summary["rotated"]++
		} else {
			summary["unchanged"]++
		}
	}
	return summary
}
