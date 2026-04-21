package envprune

import (
	"regexp"
	"strings"
)

// Result holds the outcome of pruning a single key.
type Result struct {
	Key     string
	Value   string
	Pruned  bool
	Reason  string
}

// Options controls which keys are pruned.
type Options struct {
	// RemoveEmpty removes keys with empty values.
	RemoveEmpty bool
	// RemovePrefixes removes keys matching any of these prefixes.
	RemovePrefixes []string
	// RemovePatterns removes keys matching any of these regex patterns.
	RemovePatterns []string
	// KeepOnly retains only keys matching these prefixes (if non-empty).
	KeepOnly []string
}

// Apply prunes the given env map according to opts and returns per-key results.
func Apply(env map[string]string, opts Options) []Result {
	var compiled []*regexp.Regexp
	for _, p := range opts.RemovePatterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}

	var results []Result
	for k, v := range env {
		pruned, reason := shouldPrune(k, v, opts, compiled)
		results = append(results, Result{
			Key:    k,
			Value:  v,
			Pruned: pruned,
			Reason: reason,
		})
	}
	return results
}

func shouldPrune(key, value string, opts Options, patterns []*regexp.Regexp) (bool, string) {
	if len(opts.KeepOnly) > 0 {
		for _, prefix := range opts.KeepOnly {
			if strings.HasPrefix(key, prefix) {
				return false, ""
			}
		}
		return true, "not in keep-only list"
	}
	if opts.RemoveEmpty && value == "" {
		return true, "empty value"
	}
	for _, prefix := range opts.RemovePrefixes {
		if strings.HasPrefix(key, prefix) {
			return true, "matches prefix: " + prefix
		}
	}
	for _, re := range patterns {
		if re.MatchString(key) {
			return true, "matches pattern: " + re.String()
		}
	}
	return false, ""
}

// ToMap returns only the keys that were NOT pruned.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string)
	for _, r := range results {
		if !r.Pruned {
			out[r.Key] = r.Value
		}
	}
	return out
}

// GetSummary returns counts of pruned vs retained keys.
func GetSummary(results []Result) (pruned, retained int) {
	for _, r := range results {
		if r.Pruned {
			pruned++
		} else {
			retained++
		}
	}
	return
}
