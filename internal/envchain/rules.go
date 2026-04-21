package envchain

import (
	"fmt"
	"strings"
)

// ParseStrategy parses a strategy string into a Strategy value.
func ParseStrategy(s string) (Strategy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "override", "":
		return StrategyOverride, nil
	case "keep":
		return StrategyKeep, nil
	default:
		return "", fmt.Errorf("envchain: unknown strategy %q (supported: override, keep)", s)
	}
}

// SupportedStrategies returns the list of valid strategy names.
func SupportedStrategies() []string {
	return []string{"override", "keep"}
}

// FilterBySource returns only results resolved by the given source name.
func FilterBySource(results []Result, source string) []Result {
	var out []Result
	for _, r := range results {
		if r.ResolvedBy == source {
			out = append(out, r)
		}
	}
	return out
}

// FilterOverridden returns only results where a conflict was detected.
func FilterOverridden(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Overridden {
			out = append(out, r)
		}
	}
	return out
}
