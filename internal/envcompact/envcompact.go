package envcompact

import "strings"

// Strategy defines how compaction is applied.
type Strategy string

const (
	StrategyRemoveEmpty    Strategy = "remove_empty"
	StrategyRemoveDups     Strategy = "remove_dups"
	StrategyRemoveDefaults Strategy = "remove_defaults"
)

// Result holds the outcome of compacting a single key.
type Result struct {
	Key     string
	Value   string
	Removed bool
	Reason  string
}

// Options controls compaction behaviour.
type Options struct {
	Strategies []Strategy
	Defaults   map[string]string // used with remove_defaults
}

// Apply compacts env entries according to the given options.
func Apply(env map[string]string, opts Options) []Result {
	seen := map[string]bool{}
	results := make([]Result, 0, len(env))

	for k, v := range env {
		removed := false
		reason := ""

		for _, s := range opts.Strategies {
			switch s {
			case StrategyRemoveEmpty:
				if strings.TrimSpace(v) == "" {
					removed = true
					reason = "empty value"
				}
			case StrategyRemoveDups:
				if seen[v] {
					removed = true
					reason = "duplicate value"
				}
				seen[v] = true
			case StrategyRemoveDefaults:
				if dv, ok := opts.Defaults[k]; ok && dv == v {
					removed = true
					reason = "matches default"
				}
			}
			if removed {
				break
			}
		}

		results = append(results, Result{
			Key:     k,
			Value:   v,
			Removed: removed,
			Reason:  reason,
		})
	}
	return results
}

// ToMap returns only the keys that were NOT removed.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string)
	for _, r := range results {
		if !r.Removed {
			out[r.Key] = r.Value
		}
	}
	return out
}

// GetSummary returns counts of kept and removed entries.
func GetSummary(results []Result) (kept, removed int) {
	for _, r := range results {
		if r.Removed {
			removed++
		} else {
			kept++
		}
	}
	return
}
