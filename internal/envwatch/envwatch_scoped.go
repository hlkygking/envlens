package envwatch

import "sort"

// ScopedResult holds watch results for a named scope.
type ScopedResult struct {
	Scope   string
	Entries []Entry
}

// ScopedApply runs Apply across multiple named scope maps against a shared baseline.
func ScopedApply(baseline map[string]string, scopes map[string]map[string]string, opts Options) []ScopedResult {
	names := make([]string, 0, len(scopes))
	for name := range scopes {
		names = append(names, name)
	}
	sort.Strings(names)

	results := make([]ScopedResult, 0, len(names))
	for _, name := range names {
		entries := Apply(baseline, scopes[name], opts)
		results = append(results, ScopedResult{
			Scope:   name,
			Entries: entries,
		})
	}
	return results
}

// FilterByStatus returns only entries matching the given status from a result slice.
func FilterByStatus(entries []Entry, status Status) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Status == status {
			out = append(out, e)
		}
	}
	return out
}

// HasChanges returns true if any watched entry has a non-unchanged status.
func HasChanges(entries []Entry) bool {
	for _, e := range entries {
		if e.Watched && e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}
