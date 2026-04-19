package trim

import (
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

func TestApply_TrimsWhitespace(t *testing.T) {
	env := map[string]string{
		"FOO": "  hello  ",
		"BAR": "world",
	}
	results := Apply(env, false)
	r, ok := findResult(results, "FOO")
	if !ok {
		t.Fatal("expected FOO in results")
	}
	if r.Trimmed != "hello" {
		t.Errorf("expected 'hello', got %q", r.Trimmed)
	}
	if !r.Changed {
		t.Error("expected Changed=true for FOO")
	}
	r2, _ := findResult(results, "BAR")
	if r2.Changed {
		t.Error("expected Changed=false for BAR")
	}
}

func TestApply_StripQuotes(t *testing.T) {
	env := map[string]string{
		"QUOTED": `"my value"`,
		"SINGLE": `'another'`,
		"PLAIN":  "plain",
	}
	results := Apply(env, true)
	for _, tc := range []struct{ key, want string }{
		{"QUOTED", "my value"},
		{"SINGLE", "another"},
		{"PLAIN", "plain"},
	} {
		r, ok := findResult(results, tc.key)
		if !ok {
			t.Fatalf("missing key %s", tc.key)
		}
		if r.Trimmed != tc.want {
			t.Errorf("%s: expected %q got %q", tc.key, tc.want, r.Trimmed)
		}
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"A": " val ", "B": "ok"}
	results := Apply(env, false)
	m := ToMap(results)
	if m["A"] != "val" {
		t.Errorf("expected 'val', got %q", m["A"])
	}
	if m["B"] != "ok" {
		t.Errorf("expected 'ok', got %q", m["B"])
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"X": " spaces ",
		"Y": "clean",
		"Z": "\t tab",
	}
	results := Apply(env, false)
	s := GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Changed != 2 {
		t.Errorf("expected Changed=2, got %d", s.Changed)
	}
}
