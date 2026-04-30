package envreplace

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

func TestApply_ExactMatch(t *testing.T) {
	env := map[string]string{"ENV": "staging", "DB": "localhost"}
	rules := []Rule{{Target: "staging", Replacement: "production", Strategy: StrategyExact}}
	results := Apply(env, rules)
	r, ok := findResult(results, "ENV")
	if !ok {
		t.Fatal("expected result for ENV")
	}
	if r.Status != StatusReplaced {
		t.Errorf("expected replaced, got %s", r.Status)
	}
	if r.NewValue != "production" {
		t.Errorf("expected production, got %s", r.NewValue)
	}
}

func TestApply_Unchanged(t *testing.T) {
	env := map[string]string{"ENV": "production"}
	rules := []Rule{{Target: "staging", Replacement: "production", Strategy: StrategyExact}}
	results := Apply(env, rules)
	r, ok := findResult(results, "ENV")
	if !ok {
		t.Fatal("expected result for ENV")
	}
	if r.Status != StatusUnchanged {
		t.Errorf("expected unchanged, got %s", r.Status)
	}
	if r.NewValue != "production" {
		t.Errorf("expected production, got %s", r.NewValue)
	}
}

func TestApply_PrefixStrategy(t *testing.T) {
	env := map[string]string{"HOST": "dev.example.com"}
	rules := []Rule{{Target: "dev.", Replacement: "prod.", Strategy: StrategyPrefix}}
	results := Apply(env, rules)
	r, ok := findResult(results, "HOST")
	if !ok {
		t.Fatal("expected result for HOST")
	}
	if r.Status != StatusReplaced {
		t.Errorf("expected replaced, got %s", r.Status)
	}
	if r.NewValue != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %s", r.NewValue)
	}
}

func TestApply_RegexStrategy(t *testing.T) {
	env := map[string]string{"URL": "http://localhost:8080/api"}
	rules := []Rule{{Target: `localhost:\d+`, Replacement: "prod-host", Strategy: StrategyRegex}}
	results := Apply(env, rules)
	r, ok := findResult(results, "URL")
	if !ok {
		t.Fatal("expected result for URL")
	}
	if r.Status != StatusReplaced {
		t.Errorf("expected replaced, got %s", r.Status)
	}
	if r.NewValue != "http://prod-host/api" {
		t.Errorf("unexpected value: %s", r.NewValue)
	}
}

func TestApply_InvalidRegex(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	rules := []Rule{{Target: "[invalid", Replacement: "x", Strategy: StrategyRegex}}
	results := Apply(env, rules)
	r, ok := findResult(results, "KEY")
	if !ok {
		t.Fatal("expected result for KEY")
	}
	if r.Status != StatusError {
		t.Errorf("expected error status, got %s", r.Status)
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"A": "old", "B": "keep"}
	rules := []Rule{{Target: "old", Replacement: "new", Strategy: StrategyExact}}
	results := Apply(env, rules)
	m := ToMap(results)
	if m["A"] != "new" {
		t.Errorf("expected new, got %s", m["A"])
	}
	if m["B"] != "keep" {
		t.Errorf("expected keep, got %s", m["B"])
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "staging", "B": "prod", "C": "[bad"}
	rules := []Rule{
		{Target: "staging", Replacement: "production", Strategy: StrategyExact},
		{Target: "[bad", Replacement: "x", Strategy: StrategyRegex},
	}
	results := Apply(env, rules)
	replaced, unchanged, errored := GetSummary(results)
	if replaced != 1 {
		t.Errorf("expected 1 replaced, got %d", replaced)
	}
	if unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", unchanged)
	}
	if errored != 1 {
		t.Errorf("expected 1 errored, got %d", errored)
	}
}
