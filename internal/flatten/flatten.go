package flatten

import "strings"

// Result holds a single flattened entry.
type Result struct {
	OriginalKey string
	FlatKey     string
	Value       string
	Changed     bool
}

// Summary holds aggregate stats.
type Summary struct {
	Total   int
	Changed int
}

// Options controls flattening behaviour.
type Options struct {
	Separator   string // default "_"
	Uppercase   bool
	StripPrefix string
}

// Apply flattens keys by normalising separators, optionally uppercasing and
// stripping a common prefix.
func Apply(env map[string]string, opts Options) []Result {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	results := make([]Result, 0, len(env))
	for k, v := range env {
		flat := k
		// Replace dots and dashes with the chosen separator.
		flat = strings.ReplaceAll(flat, ".", opts.Separator)
		flat = strings.ReplaceAll(flat, "-", opts.Separator)
		if opts.StripPrefix != "" {
			flat = strings.TrimPrefix(flat, opts.StripPrefix)
		}
		if opts.Uppercase {
			flat = strings.ToUpper(flat)
		}
		results = append(results, Result{
			OriginalKey: k,
			FlatKey:     flat,
			Value:       v,
			Changed:     flat != k,
		})
	}
	return results
}

// ToMap converts results to a plain map using the flat key.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.FlatKey] = r.Value
	}
	return m
}

// GetSummary returns aggregate counts.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Changed {
			s.Changed++
		}
	}
	return s
}
