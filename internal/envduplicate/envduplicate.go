package envduplicate

import "strings"

// Status represents whether a key's value is a duplicate.
type Status string

const (
	StatusUnique    Status = "unique"
	StatusDuplicate Status = "duplicate"
)

// Result holds the analysis of a single env entry.
type Result struct {
	Key        string
	Value      string
	Status     Status
	SharedWith []string // other keys that share the same value
}

// Summary holds aggregate counts.
type Summary struct {
	Total      int
	Unique     int
	Duplicates int
	Groups     int // number of distinct duplicate value groups
}

// Apply scans the provided map and identifies keys that share identical values.
func Apply(env map[string]string) []Result {
	// Build an inverted index: value -> []key
	valueIndex := make(map[string][]string)
	for k, v := range env {
		valueIndex[v] = append(valueIndex[v], k)
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sortStrings(keys)

	results := make([]Result, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		peers := valueIndex[v]
		shared := make([]string, 0, len(peers)-1)
		for _, p := range peers {
			if p != k {
				shared = append(shared, p)
			}
		}
		sortStrings(shared)

		status := StatusUnique
		if len(shared) > 0 {
			status = StatusDuplicate
		}
		results = append(results, Result{
			Key:        k,
			Value:      v,
			Status:     status,
			SharedWith: shared,
		})
	}
	return results
}

// GetSummary returns aggregate statistics for the results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	groupSeen := make(map[string]bool)
	for _, r := range results {
		if r.Status == StatusUnique {
			s.Unique++
		} else {
			s.Duplicates++
			// Represent group by sorted joined peers + self
			all := append([]string{r.Key}, r.SharedWith...)
			sortStrings(all)
			groupKey := strings.Join(all, "|")
			if !groupSeen[groupKey] {
				groupSeen[groupKey] = true
				s.Groups++
			}
		}
	}
	return s
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
