package envdoc

import (
	"fmt"
	"strings"
)

// RuleFormat describes the colon-separated format for doc entries.
const RuleFormat = "KEY:description:example:required:default"

// ParseInline parses a slice of raw rule strings into Entry values.
// Format: KEY:description:example:required:default
// Each field is optional except for KEY. The required field must be
// the literal string "true" to be treated as required.
func ParseInline(rules []string) ([]Entry, error) {
	var entries []Entry
	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 5)
		if len(parts) < 1 || strings.TrimSpace(parts[0]) == "" {
			return nil, fmt.Errorf("invalid doc rule (empty key): %q", rule)
		}
		e := Entry{Key: strings.TrimSpace(parts[0])}
		if len(parts) > 1 {
			e.Description = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 {
			e.Example = strings.TrimSpace(parts[2])
		}
		if len(parts) > 3 {
			e.Required = strings.TrimSpace(parts[3]) == "true"
		}
		if len(parts) > 4 {
			e.Default = strings.TrimSpace(parts[4])
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// FilterRequired returns only entries marked as required.
func FilterRequired(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Required {
			out = append(out, e)
		}
	}
	return out
}

// FilterUndocumented returns results that have no matching doc entry.
func FilterUndocumented(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if !r.Found {
			out = append(out, r)
		}
	}
	return out
}

// FilterDocumented returns results that have a matching doc entry.
func FilterDocumented(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Found {
			out = append(out, r)
		}
	}
	return out
}

// IndexEntries returns a map of entry key to Entry for fast lookup.
func IndexEntries(entries []Entry) map[string]Entry {
	idx := make(map[string]Entry, len(entries))
	for _, e := range entries {
		idx[e.Key] = e
	}
	return idx
}
