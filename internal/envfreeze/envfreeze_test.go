package envfreeze_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/envfreeze"
)

func findEntry(entries []envfreeze.Entry, key string) (envfreeze.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return envfreeze.Entry{}, false
}

func TestApply_Frozen(t *testing.T) {
	frozen := map[string]string{"APP_ENV": "production"}
	current := map[string]string{"APP_ENV": "production"}
	entries := envfreeze.Apply(frozen, current)
	e, ok := findEntry(entries, "APP_ENV")
	if !ok {
		t.Fatal("expected APP_ENV entry")
	}
	if e.Status != envfreeze.StatusFrozen {
		t.Errorf("expected frozen, got %s", e.Status)
	}
}

func TestApply_Drifted_ValueChanged(t *testing.T) {
	frozen := map[string]string{"DB_HOST": "prod-db"}
	current := map[string]string{"DB_HOST": "staging-db"}
	entries := envfreeze.Apply(frozen, current)
	e, ok := findEntry(entries, "DB_HOST")
	if !ok {
		t.Fatal("expected DB_HOST entry")
	}
	if e.Status != envfreeze.StatusDrifted {
		t.Errorf("expected drifted, got %s", e.Status)
	}
	if e.Frozen != "prod-db" || e.Current != "staging-db" {
		t.Errorf("unexpected values: frozen=%s current=%s", e.Frozen, e.Current)
	}
}

func TestApply_Drifted_KeyMissing(t *testing.T) {
	frozen := map[string]string{"SECRET_KEY": "abc"}
	current := map[string]string{}
	entries := envfreeze.Apply(frozen, current)
	e, ok := findEntry(entries, "SECRET_KEY")
	if !ok {
		t.Fatal("expected SECRET_KEY entry")
	}
	if e.Status != envfreeze.StatusDrifted {
		t.Errorf("expected drifted, got %s", e.Status)
	}
}

func TestApply_New(t *testing.T) {
	frozen := map[string]string{}
	current := map[string]string{"NEW_VAR": "hello"}
	entries := envfreeze.Apply(frozen, current)
	e, ok := findEntry(entries, "NEW_VAR")
	if !ok {
		t.Fatal("expected NEW_VAR entry")
	}
	if e.Status != envfreeze.StatusNew {
		t.Errorf("expected new, got %s", e.Status)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	frozen := map[string]string{"A": "1", "B": "2"}
	current := map[string]string{"A": "1", "B": "changed", "C": "3"}
	entries := envfreeze.Apply(frozen, current)
	s := envfreeze.GetSummary(entries)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Frozen != 1 {
		t.Errorf("expected 1 frozen, got %d", s.Frozen)
	}
	if s.Drifted != 1 {
		t.Errorf("expected 1 drifted, got %d", s.Drifted)
	}
	if s.New != 1 {
		t.Errorf("expected 1 new, got %d", s.New)
	}
}

func TestApply_SortedKeys(t *testing.T) {
	frozen := map[string]string{"Z": "1", "A": "1", "M": "1"}
	current := map[string]string{"Z": "1", "A": "1", "M": "1"}
	entries := envfreeze.Apply(frozen, current)
	if entries[0].Key != "A" || entries[1].Key != "M" || entries[2].Key != "Z" {
		t.Errorf("entries not sorted: %v", entries)
	}
}
