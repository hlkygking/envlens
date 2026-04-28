package envdeprecate

// DefaultRules returns a set of commonly deprecated environment variable patterns.
func DefaultRules() []Rule {
	return []Rule{
		{
			Pattern: "^OLD_",
			Reason:  "keys prefixed with OLD_ are considered legacy",
		},
		{
			Pattern: "^LEGACY_",
			Reason:  "keys prefixed with LEGACY_ are no longer supported",
		},
		{
			Pattern: "_PASS$",
			Replacement: "",
			Reason:  "use _PASSWORD suffix for clarity",
		},
		{
			Pattern: "^DEPRECATED_",
			Reason:  "explicitly marked as deprecated",
		},
	}
}

// MergeRules combines default rules with user-provided rules, user rules taking precedence.
func MergeRules(defaults, custom []Rule) []Rule {
	seen := make(map[string]bool)
	var merged []Rule
	for _, r := range custom {
		seen[r.Pattern] = true
		merged = append(merged, r)
	}
	for _, r := range defaults {
		if !seen[r.Pattern] {
			merged = append(merged, r)
		}
	}
	return merged
}

// FilterByStatus returns only results matching the given status.
func FilterByStatus(results []Result, status Status) []Result {
	var out []Result
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
