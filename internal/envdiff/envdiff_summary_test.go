package envdiff

import (
	"testing"
)

func makeEntries() []Entry {
	return []Entry{
		{Key: "A", Status: StatusAdded, NewValue: "1"},
		{Key: "B", Status: StatusRemoved, OldValue: "2"},
		{Key: "C", Status: StatusModified, OldValue: "3", NewValue: "4"},
		{Key: "D", Status: StatusUnchanged, OldValue: "5", NewValue: "5"},
		{Key: "E", Status: StatusAdded, NewValue: "6"},
	}
}

func TestCountByStatus_Counts(t *testing.T) {
	entries := makeEntries()
	s := CountByStatus(entries)

	if s.Added != 2 {
		t.Errorf("expected Added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", s.Unchanged)
	}
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
}

func TestHasChanges_True(t *testing.T) {
	s := CountByStatus(makeEntries())
	if !s.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestHasChanges_False(t *testing.T) {
	entries := []Entry{
		{Key: "X", Status: StatusUnchanged, OldValue: "v", NewValue: "v"},
	}
	s := CountByStatus(entries)
	if s.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestChangeCount(t *testing.T) {
	s := CountByStatus(makeEntries())
	if s.ChangeCount() != 4 {
		t.Errorf("expected ChangeCount=4, got %d", s.ChangeCount())
	}
}

func TestFilterByStatus_Added(t *testing.T) {
	result := FilterByStatus(makeEntries(), StatusAdded)
	if len(result) != 2 {
		t.Errorf("expected 2 added entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Status != StatusAdded {
			t.Errorf("unexpected status %q", e.Status)
		}
	}
}

func TestFilterByStatus_NoFilter(t *testing.T) {
	result := FilterByStatus(makeEntries())
	if len(result) != len(makeEntries()) {
		t.Errorf("expected all entries returned, got %d", len(result))
	}
}

func TestFilterByStatus_MultipleStatuses(t *testing.T) {
	result := FilterByStatus(makeEntries(), StatusAdded, StatusRemoved)
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}
