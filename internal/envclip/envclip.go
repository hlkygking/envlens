package envclip

import (
	"fmt"
	"strings"
)

// Result holds the outcome of clipping a single env entry.
type Result struct {
	Key      string
	Original string
	Clipped  string
	Truncated bool
}

// Options controls how values are clipped.
type Options struct {
	MaxLen    int
	Suffix    string // appended when truncated, e.g. "..."
	SkipKeys  []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxLen: 64,
		Suffix: "...",
	}
}

// Apply clips all env values that exceed MaxLen.
func Apply(env map[string]string, opts Options) []Result {
	if opts.MaxLen <= 0 {
		opts.MaxLen = DefaultOptions().MaxLen
	}

	skip := make(map[string]bool, len(opts.SkipKeys))
	for _, k := range opts.SkipKeys {
		skip[strings.ToUpper(k)] = true
	}

	results := make([]Result, 0, len(env))
	for k, v := range env {
		r := Result{Key: k, Original: v, Clipped: v}
		if !skip[strings.ToUpper(k)] && len(v) > opts.MaxLen {
			r.Clipped = v[:opts.MaxLen] + opts.Suffix
			r.Truncated = true
		}
		results = append(results, r)
	}
	return results
}

// ToMap converts results back to a plain map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Clipped
	}
	return m
}

// GetSummary returns a human-readable summary line.
func GetSummary(results []Result) string {
	total := len(results)
	truncated := 0
	for _, r := range results {
		if r.Truncated {
			truncated++
		}
	}
	return fmt.Sprintf("%d keys processed, %d truncated", total, truncated)
}
