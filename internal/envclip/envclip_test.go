package envclip

import (
	"strings"
	"testing"
)

func findResult(results []Result, key string) (Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

func TestApply_ShortValueUnchanged(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	results := Apply(env, DefaultOptions())
	r, ok := findResult(results, "HOST")
	if !ok {
		t.Fatal("expected HOST in results")
	}
	if r.Truncated {
		t.Error("short value should not be truncated")
	}
	if r.Clipped != "localhost" {
		t.Errorf("expected 'localhost', got %q", r.Clipped)
	}
}

func TestApply_LongValueTruncated(t *testing.T) {
	long := strings.Repeat("x", 100)
	env := map[string]string{"BIG": long}
	opts := Options{MaxLen: 10, Suffix: "..."}
	results := Apply(env, opts)
	r, ok := findResult(results, "BIG")
	if !ok {
		t.Fatal("expected BIG in results")
	}
	if !r.Truncated {
		t.Error("long value should be truncated")
	}
	if r.Clipped != "xxxxxxxxxx..." {
		t.Errorf("unexpected clipped value: %q", r.Clipped)
	}
}

func TestApply_SkipKeys(t *testing.T) {
	long := strings.Repeat("a", 100)
	env := map[string]string{"SECRET": long, "OTHER": long}
	opts := Options{MaxLen: 10, Suffix: "...", SkipKeys: []string{"SECRET"}}
	results := Apply(env, opts)

	sec, _ := findResult(results, "SECRET")
	if sec.Truncated {
		t.Error("SECRET should be skipped and not truncated")
	}
	other, _ := findResult(results, "OTHER")
	if !other.Truncated {
		t.Error("OTHER should be truncated")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"A": "hello", "B": strings.Repeat("z", 80)}
	opts := Options{MaxLen: 20, Suffix: "~"}
	results := Apply(env, opts)
	m := ToMap(results)
	if _, ok := m["A"]; !ok {
		t.Error("expected key A in map")
	}
	if _, ok := m["B"]; !ok {
		t.Error("expected key B in map")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Truncated: false},
		{Key: "B", Truncated: true},
		{Key: "C", Truncated: true},
	}
	s := GetSummary(results)
	if !strings.Contains(s, "3 keys") {
		t.Errorf("expected '3 keys' in summary, got: %s", s)
	}
	if !strings.Contains(s, "2 truncated") {
		t.Errorf("expected '2 truncated' in summary, got: %s", s)
	}
}
