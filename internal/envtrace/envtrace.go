package envtrace

import (
	"fmt"
	"strings"
)

// Source represents where a key's value originated.
type Source string

const (
	SourceFile    Source = "file"
	SourceEnv     Source = "env"
	SourceDefault Source = "default"
	SourceUnknown Source = "unknown"
)

// TraceEntry records the origin of a single environment variable.
type TraceEntry struct {
	Key      string
	Value    string
	Source   Source
	Origin   string // file path, env name, or default label
	Override bool   // true if this key was overridden by a higher-priority source
}

// Options controls tracing behaviour.
type Options struct {
	Sources []SourceSpec
}

// SourceSpec pairs a source kind with its origin label and key-value map.
type SourceSpec struct {
	Kind   Source
	Origin string
	Data   map[string]string
}

// Apply traces each key across all provided sources in priority order
// (first source wins unless Override is enabled).
func Apply(opts Options) []TraceEntry {
	seen := map[string]*TraceEntry{}
	var order []string

	for _, spec := range opts.Sources {
		for k, v := range spec.Data {
			if existing, ok := seen[k]; ok {
				if !existing.Override {
					existing.Override = true
				}
				continue
			}
			entry := &TraceEntry{
				Key:    k,
				Value:  v,
				Source: spec.Kind,
				Origin: spec.Origin,
			}
			seen[k] = entry
			order = append(order, k)
		}
	}

	results := make([]TraceEntry, 0, len(order))
	for _, k := range order {
		results = append(results, *seen[k])
	}
	return results
}

// GetSummary returns a human-readable summary of trace results.
func GetSummary(entries []TraceEntry) string {
	counts := map[Source]int{}
	overridden := 0
	for _, e := range entries {
		counts[e.Source]++
		if e.Override {
			overridden++
		}
	}
	parts := []string{fmt.Sprintf("total=%d", len(entries))}
	for _, src := range []Source{SourceFile, SourceEnv, SourceDefault, SourceUnknown} {
		if n := counts[src]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", src, n))
		}
	}
	if overridden > 0 {
		parts = append(parts, fmt.Sprintf("overridden=%d", overridden))
	}
	return strings.Join(parts, " ")
}

// FilterBySource returns only entries matching the given source kind.
func FilterBySource(entries []TraceEntry, src Source) []TraceEntry {
	var out []TraceEntry
	for _, e := range entries {
		if e.Source == src {
			out = append(out, e)
		}
	}
	return out
}
