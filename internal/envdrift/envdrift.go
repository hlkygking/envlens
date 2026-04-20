package envdrift

import (
	"fmt"
	"strings"
)

// Status represents the drift state of a key across environments.
type Status string

const (
	StatusMatch   Status = "match"
	StatusDrifted Status = "drifted"
	StatusMissing Status = "missing"
)

// Result holds the drift analysis for a single key.
type Result struct {
	Key      string
	Status   Status
	Values   map[string]string // env name -> value
	Baseline string
}

// Summary holds aggregate drift statistics.
type Summary struct {
	Total   int
	Match   int
	Drifted int
	Missing int
}

// Apply compares a baseline env map against one or more target env maps.
// It returns a Result for every key found in any of the inputs.
func Apply(baseline map[string]string, targets map[string]map[string]string) []Result {
	keys := allKeys(baseline, targets)
	results := make([]Result, 0, len(keys))

	for _, key := range keys {
		baseVal, inBase := baseline[key]
		values := map[string]string{"baseline": baseVal}
		hasDrift := false
		hasMissing := !inBase

		for envName, envMap := range targets {
			v, ok := envMap[key]
			if !ok {
				hasMissing = true
				values[envName] = ""
			} else {
				values[envName] = v
				if inBase && v != baseVal {
					hasDrift = true
				}
			}
		}

		var status Status
		switch {
		case hasMissing:
			status = StatusMissing
		case hasDrift:
			status = StatusDrifted
		default:
			status = StatusMatch
		}

		results = append(results, Result{
			Key:      key,
			Status:   status,
			Values:   values,
			Baseline: baseVal,
		})
	}
	return results
}

// GetSummary returns aggregate counts from a slice of Results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			s.Match++
		case StatusDrifted:
			s.Drifted++
		case StatusMissing:
			s.Missing++
		}
	}
	return s
}

// Format returns a human-readable string for a single result.
func Format(r Result) string {
	parts := []string{fmt.Sprintf("[%s] %s", strings.ToUpper(string(r.Status)), r.Key)}
	for env, val := range r.Values {
		parts = append(parts, fmt.Sprintf("  %s=%q", env, val))
	}
	return strings.Join(parts, "\n")
}

func allKeys(baseline map[string]string, targets map[string]map[string]string) []string {
	seen := map[string]struct{}{}
	for k := range baseline {
		seen[k] = struct{}{}
	}
	for _, m := range targets {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
