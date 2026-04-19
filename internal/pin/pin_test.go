package pin

import (
	"testing"
)

func findEntry(entries []Entry, key string) (Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return Entry{}, false
}

func TestCheck_AllMatch(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	pins := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	r := Check(env, pins)
	if r.Matched != 2 || r.Mismatch != 0 || r.Missing != 0 {
		t.Errorf("expected 2 matched, got %+v", r)
	}
}

func TestCheck_Mismatch(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging"}
	pins := map[string]string{"APP_ENV": "production"}
	r := Check(env, pins)
	if r.Mismatch != 1 {
		t.Errorf("expected 1 mismatch, got %d", r.Mismatch)
	}
	e, ok := findEntry(r.Entries, "APP_ENV")
	if !ok || e.Matched {
		t.Error("expected APP_ENV entry with Matched=false")
	}
	if e.Actual != "staging" || e.Expected != "production" {
		t.Errorf("unexpected entry values: %+v", e)
	}
}

func TestCheck_Missing(t *testing.T) {
	env := map[string]string{}
	pins := map[string]string{"DB_HOST": "localhost"}
	r := Check(env, pins)
	if r.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", r.Missing)
	}
	e, ok := findEntry(r.Entries, "DB_HOST")
	if !ok || e.Matched {
		t.Error("expected DB_HOST entry with Matched=false")
	}
}

func TestCheck_Mixed(t *testing.T) {
	env := map[string]string{"A": "1", "B": "wrong"}
	pins := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := Check(env, pins)
	if r.Matched != 1 || r.Mismatch != 1 || r.Missing != 1 {
		t.Errorf("unexpected counts: %+v", r)
	}
}

func TestSummary(t *testing.T) {
	r := Result{Matched: 3, Mismatch: 1, Missing: 2}
	s := Summary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
