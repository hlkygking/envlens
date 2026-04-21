package envpivot

import "sort"

// ScopeNames extracts and sorts the scope names from a scopes map.
func ScopeNames(scopes map[string]map[string]string) []string {
	names := make([]string, 0, len(scopes))
	for name := range scopes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// FilterDivergent returns only entries that differ across scopes.
func FilterDivergent(entries []Entry) []Entry {
	out := make([]Entry, 0)
	for _, e := range entries {
		if !e.Uniform {
			out = append(out, e)
		}
	}
	return out
}

// FilterMissing returns entries that are absent in at least one scope.
func FilterMissing(entries []Entry) []Entry {
	out := make([]Entry, 0)
	for _, e := range entries {
		if len(e.Missing) > 0 {
			out = append(out, e)
		}
	}
	return out
}

// ToScopeMap reconstructs a scope-keyed map from pivot entries for a given scope.
func ToScopeMap(entries []Entry, scope string) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if val, ok := e.Values[scope]; ok {
			m[e.Key] = val
		}
	}
	return m
}
