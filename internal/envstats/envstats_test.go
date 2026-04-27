package envstats_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envstats"
)

func findResult(results []envstats.Result, key string) (envstats.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return envstats.Result{}, false
}

func TestApply_BasicStats(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	results := envstats.Apply(env)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	r, ok := findResult(results, "PORT")
	if !ok {
		t.Fatal("expected PORT in results")
	}
	if r.Length != 4 {
		t.Errorf("expected length 4, got %d", r.Length)
	}
	if !r.HasDigits {
		t.Error("expected HasDigits=true for PORT")
	}
}

func TestApply_EmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY_VAR": ""}
	results := envstats.Apply(env)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].IsEmpty {
		t.Error("expected IsEmpty=true")
	}
	if results[0].Length != 0 {
		t.Errorf("expected length 0, got %d", results[0].Length)
	}
}

func TestApply_CharacterFlags(t *testing.T) {
	env := map[string]string{"SECRET": "P@ssw0rd!"}
	results := envstats.Apply(env)
	r := results[0]
	if !r.HasUpper {
		t.Error("expected HasUpper=true")
	}
	if !r.HasLower {
		t.Error("expected HasLower=true")
	}
	if !r.HasDigits {
		t.Error("expected HasDigits=true")
	}
	if !r.HasSpecial {
		t.Error("expected HasSpecial=true")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"A": "",
		"B": "hello",
		"C": "hello world",
	}
	results := envstats.Apply(env)
	s := envstats.GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
	if s.MaxLength != 11 {
		t.Errorf("expected MaxLength=11, got %d", s.MaxLength)
	}
	if s.LongestKey != "C" {
		t.Errorf("expected LongestKey=C, got %s", s.LongestKey)
	}
}

func TestGetSummary_Empty(t *testing.T) {
	s := envstats.GetSummary(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0 for empty input")
	}
}
