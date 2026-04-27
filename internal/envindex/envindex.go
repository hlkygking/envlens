package envindex

import "sort"

// Entry represents an indexed environment variable with its position and metadata.
type Entry struct {
	Key      string
	Value    string
	Index    int
	Length   int
	IsEmpty  bool
}

// Result holds the indexed entries and summary.
type Result struct {
	Entries []Entry
	Summary Summary
}

// Summary contains aggregate counts for the index.
type Summary struct {
	Total   int
	Empty   int
	NonEmpty int
}

// Apply indexes all keys in the provided map, assigning a stable
// alphabetical position to each entry and computing basic metadata.
func Apply(env map[string]string) Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for i, k := range keys {
		v := env[k]
		entries = append(entries, Entry{
			Key:    k,
			Value:  v,
			Index:  i,
			Length: len(v),
			IsEmpty: v == "",
		})
	}

	return Result{
		Entries: entries,
		Summary: getSummary(entries),
	}
}

// ToMap converts indexed entries back to a plain map.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// GetSummary returns aggregate counts from a slice of entries.
func GetSummary(entries []Entry) Summary {
	return getSummary(entries)
}

func getSummary(entries []Entry) Summary {
	s := Summary{Total: len(entries)}
	for _, e := range entries {
		if e.IsEmpty {
			s.Empty++
		} else {
			s.NonEmpty++
		}
	}
	return s
}
