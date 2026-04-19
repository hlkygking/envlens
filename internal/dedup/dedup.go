package dedup

import "strings"

// Strategy defines how to handle duplicate keys across sources.
type Strategy string

const (
	StrategyFirst Strategy = "first"
	StrategyLast  Strategy = "last"
)

// Result holds the outcome for a single key.
type Result struct {
	Key      string
	Value    string
	Sources  []string
	Kept     string // which source was kept
	Duplicate bool
}

// Apply deduplicates keys from multiple named sources.
// sources is a slice of (name, map) pairs in order.
func Apply(sources []NamedSource, strategy Strategy) []Result {
	type entry struct {
		value  string
		source string
	}
	seen := map[string][]entry{}
	order := []string{}

	for _, ns := range sources {
		for k, v := range ns.Env {
			if _, exists := seen[k]; !exists {
				order = append(order, k)
			}
			seen[k] = append(seen[k], entry{value: v, source: ns.Name})
		}
	}

	results := make([]Result, 0, len(order))
	for _, k := range order {
		entries := seen[k]
		var kept entry
		if strategy == StrategyLast {
			kept = entries[len(entries)-1]
		} else {
			kept = entries[0]
		}
		srcs := make([]string, len(entries))
		for i, e := range entries {
			srcs[i] = e.source
		}
		results = append(results, Result{
			Key:       k,
			Value:     kept.value,
			Sources:   srcs,
			Kept:      kept.source,
			Duplicate: len(entries) > 1,
		})
	}
	return results
}

// NamedSource pairs a name with an env map.
type NamedSource struct {
	Name string
	Env  map[string]string
}

// ToMap converts results to a flat map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}

// Summary returns counts.
func Summary(results []Result) (total, duplicates int) {
	for _, r := range results {
		total++
		if r.Duplicate {
			duplicates++
		}
	}
	_ = strings.ToLower // satisfy import
	return
}
