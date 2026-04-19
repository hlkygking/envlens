package envdiff

import (
	"fmt"
	"strings"
)

// ChangeType describes the kind of change for a key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single key comparison result.
type Entry struct {
	Key      string
	BaseVal  string
	TargetVal string
	Change   ChangeType
}

// Summary holds aggregate counts.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
}

// Apply compares two env maps and returns per-key diff entries.
func Apply(base, target map[string]string) []Entry {
	seen := map[string]bool{}
	var results []Entry

	for k, bv := range base {
		seen[k] = true
		if tv, ok := target[k]; !ok {
			results = append(results, Entry{Key: k, BaseVal: bv, Change: Removed})
		} else if bv != tv {
			results = append(results, Entry{Key: k, BaseVal: bv, TargetVal: tv, Change: Modified})
		} else {
			results = append(results, Entry{Key: k, BaseVal: bv, TargetVal: tv, Change: Unchanged})
		}
	}

	for k, tv := range target {
		if !seen[k] {
			results = append(results, Entry{Key: k, TargetVal: tv, Change: Added})
		}
	}

	return results
}

// GetSummary counts entries by change type.
func GetSummary(entries []Entry) Summary {
	var s Summary
	for _, e := range entries {
		switch e.Change {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		case Unchanged:
			s.Unchanged++
		}
	}
	return s
}

// Format returns a human-readable line for an entry.
func Format(e Entry) string {
	switch e.Change {
	case Added:
		return fmt.Sprintf("+ %s=%s", e.Key, e.TargetVal)
	case Removed:
		return fmt.Sprintf("- %s=%s", e.Key, e.BaseVal)
	case Modified:
		return fmt.Sprintf("~ %s: %s -> %s", e.Key, e.BaseVal, e.TargetVal)
	default:
		return fmt.Sprintf("  %s=%s", e.Key, strings.TrimSpace(e.BaseVal))
	}
}
