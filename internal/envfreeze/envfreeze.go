package envfreeze

import (
	"fmt"
	"sort"
	"strings"
)

// Status represents the freeze check result for a key.
type Status string

const (
	StatusFrozen   Status = "frozen"
	StatusDrifted  Status = "drifted"
	StatusNew      Status = "new"
)

// Entry holds the result of checking one key against a frozen snapshot.
type Entry struct {
	Key      string
	Frozen   string
	Current  string
	Status   Status
}

// Summary holds aggregate counts.
type Summary struct {
	Total   int
	Frozen  int
	Drifted int
	New     int
}

// Apply compares current env map against a frozen baseline map.
// Keys in frozen but missing from current are reported as drifted.
// Keys in current but not in frozen are reported as new.
func Apply(frozen, current map[string]string) []Entry {
	seen := map[string]bool{}
	var results []Entry

	for key, frozenVal := range frozen {
		seen[key] = true
		curVal, exists := current[key]
		if !exists {
			results = append(results, Entry{Key: key, Frozen: frozenVal, Current: "", Status: StatusDrifted})
			continue
		}
		if curVal != frozenVal {
			results = append(results, Entry{Key: key, Frozen: frozenVal, Current: curVal, Status: StatusDrifted})
		} else {
			results = append(results, Entry{Key: key, Frozen: frozenVal, Current: curVal, Status: StatusFrozen})
		}
	}

	for key, curVal := range current {
		if !seen[key] {
			results = append(results, Entry{Key: key, Frozen: "", Current: curVal, Status: StatusNew})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// GetSummary returns aggregate counts from a slice of entries.
func GetSummary(entries []Entry) Summary {
	s := Summary{Total: len(entries)}
	for _, e := range entries {
		switch e.Status {
		case StatusFrozen:
			s.Frozen++
		case StatusDrifted:
			s.Drifted++
		case StatusNew:
			s.New++
		}
	}
	return s
}

// Format returns a human-readable line for an entry.
func Format(e Entry) string {
	switch e.Status {
	case StatusFrozen:
		return fmt.Sprintf("[frozen]  %s", e.Key)
	case StatusDrifted:
		return fmt.Sprintf("[drifted] %s: %q -> %q", e.Key, e.Frozen, e.Current)
	case StatusNew:
		return fmt.Sprintf("[new]     %s=%q", e.Key, e.Current)
	}
	return strings.ToLower(string(e.Status))
}
