package envdiff

import (
	"testing"
)

func findEntry(entries []Entry, key string) *Entry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}

func TestApply_Added(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "1", "B": "2"}
	entries := Apply(base, target)
	e := findEntry(entries, "B")
	if e == nil || e.Change != Added {
		t.Fatalf("expected B to be added")
	}
}

func TestApply_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1"}
	entries := Apply(base, target)
	e := findEntry(entries, "B")
	if e == nil || e.Change != Removed {
		t.Fatalf("expected B to be removed")
	}
}

func TestApply_Modified(t *testing.T) {
	base := map[string]string{"A": "old"}
	target := map[string]string{"A": "new"}
	entries := Apply(base, target)
	e := findEntry(entries, "A")
	if e == nil || e.Change != Modified {
		t.Fatalf("expected A to be modified")
	}
	if e.BaseVal != "old" || e.TargetVal != "new" {
		t.Errorf("unexpected values: %+v", e)
	}
}

func TestApply_Unchanged(t *testing.T) {
	base := map[string]string{"A": "same"}
	target := map[string]string{"A": "same"}
	entries := Apply(base, target)
	e := findEntry(entries, "A")
	if e == nil || e.Change != Unchanged {
		t.Fatalf("expected A to be unchanged")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	entries := []Entry{
		{Key: "A", Change: Added},
		{Key: "B", Change: Removed},
		{Key: "C", Change: Modified},
		{Key: "D", Change: Unchanged},
		{Key: "E", Change: Unchanged},
	}
	s := GetSummary(entries)
	if s.Added != 1 || s.Removed != 1 || s.Modified != 1 || s.Unchanged != 2 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestFormat_Lines(t *testing.T) {
	if got := Format(Entry{Key: "X", TargetVal: "1", Change: Added}); got != "+ X=1" {
		t.Errorf("unexpected: %s", got)
	}
	if got := Format(Entry{Key: "X", BaseVal: "1", Change: Removed}); got != "- X=1" {
		t.Errorf("unexpected: %s", got)
	}
}
