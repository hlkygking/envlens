package envchain

import "fmt"

// Strategy controls how chained envs are merged.
type Strategy string

const (
	StrategyOverride Strategy = "override" // later sources win
	StrategyKeep     Strategy = "keep"     // first source wins
)

// Source represents a named env map in the chain.
type Source struct {
	Name string
	Env  map[string]string
}

// Result holds the resolved value for a key.
type Result struct {
	Key        string
	Value      string
	ResolvedBy string // name of the source that provided the value
	Overridden bool   // true if an earlier source also had this key
}

// Apply resolves keys across chained sources according to strategy.
func Apply(sources []Source, strategy Strategy) ([]Result, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("envchain: no sources provided")
	}
	if strategy != StrategyOverride && strategy != StrategyKeep {
		return nil, fmt.Errorf("envchain: unknown strategy %q", strategy)
	}

	seen := map[string]Result{}

	for _, src := range sources {
		for k, v := range src.Env {
			existing, exists := seen[k]
			if !exists {
				seen[k] = Result{Key: k, Value: v, ResolvedBy: src.Name}
			} else if strategy == StrategyOverride {
				seen[k] = Result{Key: k, Value: v, ResolvedBy: src.Name, Overridden: true}
			} else {
				// keep: mark that a later source tried to override
				existing.Overridden = true
				seen[k] = existing
			}
		}
	}

	results := make([]Result, 0, len(seen))
	for _, r := range seen {
		results = append(results, r)
	}
	return results, nil
}

// ToMap converts results to a plain key→value map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}

// GetSummary returns counts by source.
func GetSummary(results []Result) map[string]int {
	summary := map[string]int{}
	for _, r := range results {
		summary[r.ResolvedBy]++
	}
	return summary
}
