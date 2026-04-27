package envsample_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envsample"
)

func findResult(results []envsample.Result, key string) (envsample.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return envsample.Result{}, false
}

func sampleEnv() map[string]string {
	return map[string]string{
		"ALPHA": "1",
		"BETA":  "2",
		"GAMMA": "3",
		"DELTA": "4",
		"ECHO":  "5",
	}
}

func TestApply_RandomIncludesN(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 3, Strategy: envsample.StrategyRandom, Seed: 99}
	results := envsample.Apply(env, opts)

	included := 0
	for _, r := range results {
		if r.Included {
			included++
		}
	}
	if included != 3 {
		t.Errorf("expected 3 included, got %d", included)
	}
}

func TestApply_FirstStrategy(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 2, Strategy: envsample.StrategyFirst}
	results := envsample.Apply(env, opts)

	// Alphabetically first two keys should be ALPHA and BETA.
	alpha, _ := findResult(results, "ALPHA")
	beta, _ := findResult(results, "BETA")
	if !alpha.Included {
		t.Error("expected ALPHA to be included")
	}
	if !beta.Included {
		t.Error("expected BETA to be included")
	}
}

func TestApply_LastStrategy(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 2, Strategy: envsample.StrategyLast}
	results := envsample.Apply(env, opts)

	// Alphabetically last two keys should be GAMMA and ECHO.
	gamma, _ := findResult(results, "GAMMA")
	echo, _ := findResult(results, "ECHO")
	if !gamma.Included {
		t.Error("expected GAMMA to be included")
	}
	if !echo.Included {
		t.Error("expected ECHO to be included")
	}
}

func TestApply_NLargerThanEnv(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 100, Strategy: envsample.StrategyFirst}
	results := envsample.Apply(env, opts)

	for _, r := range results {
		if !r.Included {
			t.Errorf("expected all keys included when N > len(env), but %s was excluded", r.Key)
		}
	}
}

func TestToMap_OnlyIncluded(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 2, Strategy: envsample.StrategyFirst}
	results := envsample.Apply(env, opts)
	m := envsample.ToMap(results)
	if len(m) != 2 {
		t.Errorf("expected 2 keys in map, got %d", len(m))
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := sampleEnv()
	opts := envsample.Options{N: 3, Strategy: envsample.StrategyFirst}
	results := envsample.Apply(env, opts)
	inc, exc := envsample.GetSummary(results)
	if inc != 3 {
		t.Errorf("expected 3 included, got %d", inc)
	}
	if exc != 2 {
		t.Errorf("expected 2 excluded, got %d", exc)
	}
}
