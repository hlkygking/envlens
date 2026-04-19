package dedup

import (
	"testing"
)

func findResult(results []Result, key string) *Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_NoDuplicates(t *testing.T) {
	sources := []NamedSource{
		{Name: "base", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "override", Env: map[string]string{"C": "3"}},
	}
	results := Apply(sources, StrategyFirst)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Duplicate {
			t.Errorf("key %s should not be duplicate", r.Key)
		}
	}
}

func TestApply_FirstStrategy(t *testing.T) {
	sources := []NamedSource{
		{Name: "base", Env: map[string]string{"X": "base-val"}},
		{Name: "prod", Env: map[string]string{"X": "prod-val"}},
	}
	results := Apply(sources, StrategyFirst)
	r := findResult(results, "X")
	if r == nil {
		t.Fatal("key X not found")
	}
	if r.Value != "base-val" {
		t.Errorf("expected base-val, got %s", r.Value)
	}
	if r.Kept != "base" {
		t.Errorf("expected kept=base, got %s", r.Kept)
	}
	if !r.Duplicate {
		t.Error("expected Duplicate=true")
	}
}

func TestApply_LastStrategy(t *testing.T) {
	sources := []NamedSource{
		{Name: "base", Env: map[string]string{"X": "base-val"}},
		{Name: "prod", Env: map[string]string{"X": "prod-val"}},
	}
	results := Apply(sources, StrategyLast)
	r := findResult(results, "X")
	if r == nil {
		t.Fatal("key X not found")
	}
	if r.Value != "prod-val" {
		t.Errorf("expected prod-val, got %s", r.Value)
	}
	if r.Kept != "prod" {
		t.Errorf("expected kept=prod, got %s", r.Kept)
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	sources := []NamedSource{
		{Name: "a", Env: map[string]string{"K": "v"}},
	}
	m := ToMap(Apply(sources, StrategyFirst))
	if m["K"] != "v" {
		t.Errorf("expected v, got %s", m["K"])
	}
}

func TestSummary_Counts(t *testing.T) {
	sources := []NamedSource{
		{Name: "a", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "b", Env: map[string]string{"B": "3", "C": "4"}},
	}
	results := Apply(sources, StrategyFirst)
	total, dups := Summary(results)
	if total != 3 {
		t.Errorf("expected total=3, got %d", total)
	}
	if dups != 1 {
		t.Errorf("expected dups=1, got %d", dups)
	}
}
