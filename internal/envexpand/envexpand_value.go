package envexpand

// Value is a helper method on Result to return the best available value.
// Returns Expanded if status is ok or unchanged, otherwise Original.
func (r Result) Value() string {
	if r.Status == "ok" || r.Status == "unchanged" {
		return r.Expanded
	}
	return r.Original
}

// IsExpanded returns true if the result was successfully expanded.
func (r Result) IsExpanded() bool {
	return r.Status == "ok"
}

// IsUnresolved returns true if one or more references could not be resolved.
func (r Result) IsUnresolved() bool {
	return r.Status == "unresolved"
}

// IsUnchanged returns true if the value contained no variable references.
func (r Result) IsUnchanged() bool {
	return r.Status == "unchanged"
}

// FilterByStatus returns only results matching the given status string.
func FilterByStatus(results []Result, status string) []Result {
	out := make([]Result, 0)
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}

// HasUnresolved returns true if any result has unresolved references.
func HasUnresolved(results []Result) bool {
	for _, r := range results {
		if r.IsUnresolved() {
			return true
		}
	}
	return false
}
