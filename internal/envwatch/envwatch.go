package envwatch

import (
	"fmt"
	"sort"
	"strings"
)

// Status represents the watch status of a key.
type Status string

const (
	StatusNew      Status = "new"
	StatusRemoved  Status = "removed"
	StatusChanged  Status = "changed"
	StatusUnchanged Status = "unchanged"
)

// Entry holds the watch result for a single key.
type Entry struct {
	Key       string
	OldValue  string
	NewValue  string
	Status    Status
	Watched   bool
}

// Options configures the watch behaviour.
type Options struct {
	// Keys is the explicit list of keys to watch. Empty means watch all.
	Keys []string
	// IgnoreCase makes key matching case-insensitive.
	IgnoreCase bool
}

// Apply compares baseline and current env maps and returns watch entries.
func Apply(baseline, current map[string]string, opts Options) []Entry {
	watchSet := buildWatchSet(opts.Keys, opts.IgnoreCase)

	allKeys := mergeKeys(baseline, current)
	sort.Strings(allKeys)

	var results []Entry
	for _, key := range allKeys {
		watched := len(watchSet) == 0 || isWatched(key, watchSet, opts.IgnoreCase)
		oldVal, inBase := baseline[key]
		newVal, inCurr := current[key]

		var status Status
		switch {
		case inBase && !inCurr:
			status = StatusRemoved
		case !inBase && inCurr:
			status = StatusNew
		case oldVal != newVal:
			status = StatusChanged
		default:
			status = StatusUnchanged
		}

		results = append(results, Entry{
			Key:      key,
			OldValue: oldVal,
			NewValue: newVal,
			Status:   status,
			Watched:  watched,
		})
	}
	return results
}

// GetSummary returns counts by status for watched entries.
func GetSummary(entries []Entry) map[string]int {
	summary := map[string]int{
		"new": 0, "removed": 0, "changed": 0, "unchanged": 0, "total": 0,
	}
	for _, e := range entries {
		if !e.Watched {
			continue
		}
		summary[string(e.Status)]++
		summary["total"]++
	}
	return summary
}

// Format returns a human-readable label for an entry.
func Format(e Entry) string {
	switch e.Status {
	case StatusNew:
		return fmt.Sprintf("+%s=%s", e.Key, e.NewValue)
	case StatusRemoved:
		return fmt.Sprintf("-%s=%s", e.Key, e.OldValue)
	case StatusChanged:
		return fmt.Sprintf("~%s: %s -> %s", e.Key, e.OldValue, e.NewValue)
	default:
		return fmt.Sprintf(" %s=%s", e.Key, e.NewValue)
	}
}

func buildWatchSet(keys []string, ignoreCase bool) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if ignoreCase {
			k = strings.ToUpper(k)
		}
		set[k] = struct{}{}
	}
	return set
}

func isWatched(key string, set map[string]struct{}, ignoreCase bool) bool {
	if ignoreCase {
		key = strings.ToUpper(key)
	}
	_, ok := set[key]
	return ok
}

func mergeKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
