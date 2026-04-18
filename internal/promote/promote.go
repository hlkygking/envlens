package promote

import (
	"fmt"
	"sort"
)

// Strategy controls how conflicts are resolved during promotion.
type Strategy string

const (
	StrategyOverride Strategy = "override" // target value wins
	StrategyKeep     Strategy = "keep"     // source value wins
	StrategyError    Strategy = "error"    // conflict is an error
)

// Result describes what happened to a single key during promotion.
type Result struct {
	Key      string
	Value    string
	Status   string // promoted, skipped, conflict
	Conflict bool
	Message  string
}

// Summary holds aggregate counts from a promotion run.
type Summary struct {
	Promoted  int
	Skipped   int
	Conflicts int
}

// Apply promotes keys from src into dst according to the given strategy.
// Keys already present in dst are treated as conflicts.
func Apply(src, dst map[string]string, strategy Strategy) ([]Result, error) {
	var results []Result

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		srcVal := src[k]
		dstVal, exists := dst[k]

		if !exists {
			dst[k] = srcVal
			results = append(results, Result{Key: k, Value: srcVal, Status: "promoted"})
			continue
		}

		if srcVal == dstVal {
			results = append(results, Result{Key: k, Value: srcVal, Status: "skipped", Message: "identical"})
			continue
		}

		switch strategy {
		case StrategyOverride:
			dst[k] = srcVal
			results = append(results, Result{Key: k, Value: srcVal, Status: "promoted", Conflict: true, Message: fmt.Sprintf("overrode %q", dstVal)})
		case StrategyKeep:
			results = append(results, Result{Key: k, Value: dstVal, Status: "skipped", Conflict: true, Message: "kept existing"})
		case StrategyError:
			return nil, fmt.Errorf("conflict on key %q: src=%q dst=%q", k, srcVal, dstVal)
		default:
			return nil, fmt.Errorf("unknown strategy %q", strategy)
		}
	}

	return results, nil
}

// GetSummary aggregates result counts.
func GetSummary(results []Result) Summary {
	var s Summary
	for _, r := range results {
		switch r.Status {
		case "promoted":
			s.Promoted++
			if r.Conflict {
				s.Conflicts++
			}
		case "skipped":
			s.Skipped++
		}
	}
	return s
}
