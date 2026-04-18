package diff

import "sort"

// Status represents the kind of change detected for a key.
type Status string

const (
	Added    Status = "Added"
	Removed  Status = "Removed"
	Modified Status = "Modified"
	Unchanged Status = "Unchanged"
)

// Entry describes a single key comparison result.
type Entry struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// Compare diffs two env maps (base vs target) and returns a sorted list of entries.
func Compare(base, target map[string]string) []Entry {
	var entries []Entry

	for k, bv := range base {
		if tv, ok := target[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: Removed, OldValue: bv})
		} else if bv != tv {
			entries = append(entries, Entry{Key: k, Status: Modified, OldValue: bv, NewValue: tv})
		} else {
			entries = append(entries, Entry{Key: k, Status: Unchanged, OldValue: bv, NewValue: tv})
		}
	}

	for k, tv := range target {
		if _, ok := base[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: Added, NewValue: tv})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}
