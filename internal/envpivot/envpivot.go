// Package envpivot transposes a multi-scope env map into a key-centric view,
// showing how each key's value differs across named scopes.
package envpivot

import "sort"

// Entry represents a single key's values across all scopes.
type Entry struct {
	Key      string            `json:"key"`
	Values   map[string]string `json:"values"`   // scope -> value
	Missing  []string          `json:"missing"`  // scopes where key is absent
	Uniform  bool              `json:"uniform"`  // true when all present values are equal
}

// Summary holds aggregate statistics for a pivot operation.
type Summary struct {
	TotalKeys    int `json:"total_keys"`
	UniformKeys  int `json:"uniform_keys"`
	DivergentKeys int `json:"divergent_keys"`
	MissingInAny int `json:"missing_in_any"`
}

// Apply pivots scoped env maps into a key-centric slice of Entry.
// scopes is a map of scope-name -> (key -> value).
func Apply(scopes map[string]map[string]string) []Entry {
	if len(scopes) == 0 {
		return nil
	}

	// Collect all unique keys across all scopes.
	keySet := map[string]struct{}{}
	for _, env := range scopes {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	scopeNames := make([]string, 0, len(scopes))
	for name := range scopes {
		scopeNames = append(scopeNames, name)
	}
	sort.Strings(scopeNames)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, key := range keys {
		e := Entry{
			Key:    key,
			Values: map[string]string{},
		}
		var firstVal string
		firstSet := false
		uniform := true

		for _, scope := range scopeNames {
			env := scopes[scope]
			if val, ok := env[key]; ok {
				e.Values[scope] = val
				if !firstSet {
					firstVal = val
					firstSet = true
				} else if val != firstVal {
					uniform = false
				}
			} else {
				e.Missing = append(e.Missing, scope)
				uniform = false
			}
		}
		e.Uniform = uniform
		entries = append(entries, e)
	}
	return entries
}

// GetSummary computes aggregate statistics over pivot entries.
func GetSummary(entries []Entry) Summary {
	s := Summary{TotalKeys: len(entries)}
	for _, e := range entries {
		if e.Uniform {
			s.UniformKeys++
		} else {
			s.DivergentKeys++
		}
		if len(e.Missing) > 0 {
			s.MissingInAny++
		}
	}
	return s
}
