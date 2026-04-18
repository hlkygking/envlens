package merge

import (
	"fmt"
	"sort"
)

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy string

const (
	StrategyBase     Strategy = "base"     // base wins on conflict
	StrategyOverride Strategy = "override" // override wins on conflict (default)
	StrategyError    Strategy = "error"    // conflict is an error
)

// Conflict records a key that existed in both maps with different values.
type Conflict struct {
	Key       string
	BaseValue string
	OverValue string
}

// Result holds the merged map and metadata.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
	Added     []string // keys only in override
	Kept      []string // keys only in base
}

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	return fmt.Sprintf("merged=%d conflicts=%d added=%d kept=%d",
		len(r.Merged), len(r.Conflicts), len(r.Added), len(r.Kept))
}

// Apply merges override into base using the given strategy.
func Apply(base, override map[string]string, strategy Strategy) (Result, error) {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []Conflict
	var added []string

	keys := make([]string, 0, len(override))
	for k := range override {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		ov := override[k]
		bv, exists := base[k]
		if !exists {
			merged[k] = ov
			added = append(added, k)
			continue
		}
		if bv == ov {
			continue
		}
		conflicts = append(conflicts, Conflict{Key: k, BaseValue: bv, OverValue: ov})
		switch strategy {
		case StrategyBase:
			merged[k] = bv
		case StrategyError:
			return Result{}, fmt.Errorf("conflict on key %q: base=%q override=%q", k, bv, ov)
		default: // StrategyOverride
			merged[k] = ov
		}
	}

	var kept []string
	for k := range base {
		if _, ok := override[k]; !ok {
			kept = append(kept, k)
		}
	}
	sort.Strings(kept)

	return Result{Merged: merged, Conflicts: conflicts, Added: added, Kept: kept}, nil
}
