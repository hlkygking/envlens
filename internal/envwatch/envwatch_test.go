package envwatch_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/envwatch"
)

func findResult(entries []envwatch.Entry, key string) (envwatch.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return envwatch.Entry{}, false
}

func TestApply_New(t *testing.T) {
	base := map[string]string{"A": "1"}
	curr := map[string]string{"A": "1", "B": "2"}
	entries := envwatch.Apply(base, curr, envwatch.Options{})
	e, ok := findResult(entries, "B")
	if !ok {
		t.Fatal("expected entry for B")
	}
	if e.Status != envwatch.StatusNew {
		t.Errorf("expected new, got %s", e.Status)
	}
	if e.NewValue != "2" {
		t.Errorf("expected NewValue=2, got %s", e.NewValue)
	}
}

func TestApply_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "1"}
	entries := envwatch.Apply(base, curr, envwatch.Options{})
	e, ok := findResult(entries, "B")
	if !ok {
		t.Fatal("expected entry for B")
	}
	if e.Status != envwatch.StatusRemoved {
		t.Errorf("expected removed, got %s", e.Status)
	}
	if e.OldValue != "2" {
		t.Errorf("expected OldValue=2, got %s", e.OldValue)
	}
}

func TestApply_Changed(t *testing.T) {
	base := map[string]string{"A": "old"}
	curr := map[string]string{"A": "new"}
	entries := envwatch.Apply(base, curr, envwatch.Options{})
	e, ok := findResult(entries, "A")
	if !ok {
		t.Fatal("expected entry for A")
	}
	if e.Status != envwatch.StatusChanged {
		t.Errorf("expected changed, got %s", e.Status)
	}
}

func TestApply_Unchanged(t *testing.T) {
	base := map[string]string{"A": "1"}
	curr := map[string]string{"A": "1"}
	entries := envwatch.Apply(base, curr, envwatch.Options{})
	e, ok := findResult(entries, "A")
	if !ok {
		t.Fatal("expected entry for A")
	}
	if e.Status != envwatch.StatusUnchanged {
		t.Errorf("expected unchanged, got %s", e.Status)
	}
}

func TestApply_WatchedFilter(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "99", "B": "2"}
	opts := envwatch.Options{Keys: []string{"A"}}
	entries := envwatch.Apply(base, curr, opts)

	a, _ := findResult(entries, "A")
	b, _ := findResult(entries, "B")
	if !a.Watched {
		t.Error("A should be watched")
	}
	if b.Watched {
		t.Error("B should not be watched")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "changed", "C": "3"}
	entries := envwatch.Apply(base, curr, envwatch.Options{})
	summary := envwatch.GetSummary(entries)

	if summary["changed"] != 1 {
		t.Errorf("expected 1 changed, got %d", summary["changed"])
	}
	if summary["new"] != 1 {
		t.Errorf("expected 1 new, got %d", summary["new"])
	}
	if summary["removed"] != 1 {
		t.Errorf("expected 1 removed, got %d", summary["removed"])
	}
}

func TestFormat_Labels(t *testing.T) {
	tests := []struct {
		entry  envwatch.Entry
		prefix string
	}{
		{envwatch.Entry{Key: "A", NewValue: "1", Status: envwatch.StatusNew}, "+"},
		{envwatch.Entry{Key: "B", OldValue: "2", Status: envwatch.StatusRemoved}, "-"},
		{envwatch.Entry{Key: "C", OldValue: "x", NewValue: "y", Status: envwatch.StatusChanged}, "~"},
		{envwatch.Entry{Key: "D", NewValue: "z", Status: envwatch.StatusUnchanged}, " "},
	}
	for _, tt := range tests {
		out := envwatch.Format(tt.entry)
		if len(out) == 0 || string(out[0]) != tt.prefix {
			t.Errorf("key %s: expected prefix %q, got %q", tt.entry.Key, tt.prefix, out)
		}
	}
}
