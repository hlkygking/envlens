package envindex_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envindex"
)

func findResult(entries []envindex.Entry, key string) (envindex.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return envindex.Entry{}, false
}

func TestApply_AssignsStableIndex(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"BETA":  "b",
	}
	res := envindex.Apply(env)

	a, _ := findResult(res.Entries, "ALPHA")
	b, _ := findResult(res.Entries, "BETA")
	z, _ := findResult(res.Entries, "ZEBRA")

	if a.Index != 0 {
		t.Errorf("expected ALPHA index 0, got %d", a.Index)
	}
	if b.Index != 1 {
		t.Errorf("expected BETA index 1, got %d", b.Index)
	}
	if z.Index != 2 {
		t.Errorf("expected ZEBRA index 2, got %d", z.Index)
	}
}

func TestApply_SetsLength(t *testing.T) {
	env := map[string]string{"FOO": "hello"}
	res := envindex.Apply(env)
	e, ok := findResult(res.Entries, "FOO")
	if !ok {
		t.Fatal("FOO not found")
	}
	if e.Length != 5 {
		t.Errorf("expected length 5, got %d", e.Length)
	}
}

func TestApply_DetectsEmpty(t *testing.T) {
	env := map[string]string{"EMPTY": "", "FULL": "val"}
	res := envindex.Apply(env)

	e, _ := findResult(res.Entries, "EMPTY")
	f, _ := findResult(res.Entries, "FULL")

	if !e.IsEmpty {
		t.Error("expected EMPTY to be marked empty")
	}
	if f.IsEmpty {
		t.Error("expected FULL to not be marked empty")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "", "B": "val", "C": ""}
	res := envindex.Apply(env)
	s := res.Summary

	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Empty != 2 {
		t.Errorf("expected empty 2, got %d", s.Empty)
	}
	if s.NonEmpty != 1 {
		t.Errorf("expected non-empty 1, got %d", s.NonEmpty)
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	res := envindex.Apply(env)
	out := envindex.ToMap(res.Entries)

	for k, v := range env {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}
