package envchain

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

func TestApply_Override(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"FOO": "base_foo", "BAR": "base_bar"}},
		{Name: "prod", Env: map[string]string{"FOO": "prod_foo", "BAZ": "prod_baz"}},
	}
	res, err := Apply(sources, StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := findResult(res, "FOO")
	if !ok {
		t.Fatal("expected FOO in results")
	}
	if r.Value != "prod_foo" {
		t.Errorf("expected prod_foo, got %s", r.Value)
	}
	if !r.Overridden {
		t.Error("expected Overridden=true for FOO")
	}
	if r.ResolvedBy != "prod" {
		t.Errorf("expected ResolvedBy=prod, got %s", r.ResolvedBy)
	}
}

func TestApply_Keep(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"FOO": "base_foo"}},
		{Name: "prod", Env: map[string]string{"FOO": "prod_foo"}},
	}
	res, err := Apply(sources, StrategyKeep)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := findResult(res, "FOO")
	if !ok {
		t.Fatal("expected FOO")
	}
	if r.Value != "base_foo" {
		t.Errorf("expected base_foo, got %s", r.Value)
	}
	if r.ResolvedBy != "base" {
		t.Errorf("expected ResolvedBy=base, got %s", r.ResolvedBy)
	}
}

func TestApply_NoSources(t *testing.T) {
	_, err := Apply(nil, StrategyOverride)
	if err == nil {
		t.Error("expected error for empty sources")
	}
}

func TestApply_UnknownStrategy(t *testing.T) {
	sources := []Source{{Name: "a", Env: map[string]string{"X": "1"}}}
	_, err := Apply(sources, Strategy("unknown"))
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"A": "1", "B": "2"}},
	}
	res, _ := Apply(sources, StrategyOverride)
	m := ToMap(res)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	sources := []Source{
		{Name: "base", Env: map[string]string{"A": "1", "B": "2"}},
		{Name: "prod", Env: map[string]string{"C": "3"}},
	}
	res, _ := Apply(sources, StrategyOverride)
	summary := GetSummary(res)
	if summary["base"] != 2 {
		t.Errorf("expected base=2, got %d", summary["base"])
	}
	if summary["prod"] != 1 {
		t.Errorf("expected prod=1, got %d", summary["prod"])
	}
}
