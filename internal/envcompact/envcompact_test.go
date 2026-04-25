package envcompact_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envcompact"
)

func findResult(results []envcompact.Result, key string) *envcompact.Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_RemoveEmpty(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "", "C": "  "}
	opts := envcompact.Options{Strategies: []envcompact.Strategy{envcompact.StrategyRemoveEmpty}}
	results := envcompact.Apply(env, opts)

	if r := findResult(results, "A"); r == nil || r.Removed {
		t.Error("expected A to be kept")
	}
	if r := findResult(results, "B"); r == nil || !r.Removed {
		t.Error("expected B to be removed")
	}
	if r := findResult(results, "C"); r == nil || !r.Removed {
		t.Error("expected C (whitespace) to be removed")
	}
}

func TestApply_RemoveDuplicates(t *testing.T) {
	env := map[string]string{"X": "same", "Y": "same", "Z": "unique"}
	opts := envcompact.Options{Strategies: []envcompact.Strategy{envcompact.StrategyRemoveDups}}
	results := envcompact.Apply(env, opts)

	kept := 0
	for _, r := range results {
		if !r.Removed {
			kept++
		}
	}
	// Only one of X/Y should survive; Z always kept.
	if kept != 2 {
		t.Errorf("expected 2 kept entries, got %d", kept)
	}
}

func TestApply_RemoveDefaults(t *testing.T) {
	env := map[string]string{"PORT": "8080", "HOST": "custom"}
	defaults := map[string]string{"PORT": "8080", "HOST": "localhost"}
	opts := envcompact.Options{
		Strategies: []envcompact.Strategy{envcompact.StrategyRemoveDefaults},
		Defaults:   defaults,
	}
	results := envcompact.Apply(env, opts)

	if r := findResult(results, "PORT"); r == nil || !r.Removed {
		t.Error("expected PORT to be removed (matches default)")
	}
	if r := findResult(results, "HOST"); r == nil || r.Removed {
		t.Error("expected HOST to be kept (differs from default)")
	}
}

func TestToMap_OnlyKept(t *testing.T) {
	env := map[string]string{"A": "val", "B": ""}
	opts := envcompact.Options{Strategies: []envcompact.Strategy{envcompact.StrategyRemoveEmpty}}
	m := envcompact.ToMap(envcompact.Apply(env, opts))

	if _, ok := m["B"]; ok {
		t.Error("expected B to be absent from map")
	}
	if m["A"] != "val" {
		t.Errorf("expected A=val, got %q", m["A"])
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "x", "B": "", "C": ""}
	opts := envcompact.Options{Strategies: []envcompact.Strategy{envcompact.StrategyRemoveEmpty}}
	kept, removed := envcompact.GetSummary(envcompact.Apply(env, opts))

	if kept != 1 {
		t.Errorf("expected 1 kept, got %d", kept)
	}
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
}

func TestApply_NoStrategies(t *testing.T) {
	// With no strategies applied, all entries should be kept.
	env := map[string]string{"A": "", "B": "value", "C": "  "}
	opts := envcompact.Options{Strategies: []envcompact.Strategy{}}
	results := envcompact.Apply(env, opts)

	if len(results) != len(env) {
		t.Errorf("expected %d results, got %d", len(env), len(results))
	}
	for _, r := range results {
		if r.Removed {
			t.Errorf("expected key %q to be kept with no strategies, but it was removed", r.Key)
		}
	}
}
