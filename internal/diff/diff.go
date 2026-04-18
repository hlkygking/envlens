// Package diff computes the difference between two EnvMaps.
package diff

import (
	"github.com/envlens/envlens/internal/parser"
)

// ChangeKind describes the type of change for a key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Entry represents a single diff entry.
type Entry struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if there are any non-unchanged entries.
func (r Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs base against head and returns a Result.
func Compare(base, head parser.EnvMap) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, baseVal := range base {
		seen[k] = true
		if headVal, ok := head[k]; !ok {
			entries = append(entries, Entry{Key: k, Kind: Removed, OldValue: baseVal})
		} else if baseVal != headVal {
			entries = append(entries, Entry{Key: k, Kind: Modified, OldValue: baseVal, NewValue: headVal})
		} else {
			entries = append(entries, Entry{Key: k, Kind: Unchanged, OldValue: baseVal, NewValue: headVal})
		}
	}

	for k, headVal := range head {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Kind: Added, NewValue: headVal})
		}
	}

	return Result{Entries: entries}
}
