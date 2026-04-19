package pin

import (
	"fmt"
	"sort"
)

// Entry represents a pinned key-value pair with its expected value.
type Entry struct {
	Key      string `json:"key"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Matched  bool   `json:"matched"`
}

// Result holds the outcome of a pin check.
type Result struct {
	Entries  []Entry `json:"entries"`
	Matched  int     `json:"matched"`
	Mismatch int     `json:"mismatch"`
	Missing  int     `json:"missing"`
}

// Check verifies that each pinned key in pins matches the value in env.
// pins is a map of key -> expected value.
func Check(env map[string]string, pins map[string]string) Result {
	keys := make([]string, 0, len(pins))
	for k := range pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var entries []Entry
	matched, mismatch, missing := 0, 0, 0

	for _, k := range keys {
		expected := pins[k]
		actual, ok := env[k]
		if !ok {
			entries = append(entries, Entry{Key: k, Expected: expected, Actual: "", Matched: false})
			missing++
			continue
		}
		if actual == expected {
			entries = append(entries, Entry{Key: k, Expected: expected, Actual: actual, Matched: true})
			matched++
		} else {
			entries = append(entries, Entry{Key: k, Expected: expected, Actual: actual, Matched: false})
			mismatch++
		}
	}

	return Result{Entries: entries, Matched: matched, Mismatch: mismatch, Missing: missing}
}

// Summary returns a human-readable summary line.
func Summary(r Result) string {
	return fmt.Sprintf("pinned: %d matched, %d mismatched, %d missing", r.Matched, r.Mismatch, r.Missing)
}
