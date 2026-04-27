package envrotate

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

func TestApply_SensitiveRotatedByDefault(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "abc123",
		"APP_NAME": "myapp",
	}
	results := Apply(env, DefaultOptions())
	r, ok := findResult(results, "API_KEY")
	if !ok {
		t.Fatal("expected API_KEY in results")
	}
	if !r.Rotated {
		t.Error("expected API_KEY to be rotated")
	}
	r2, _ := findResult(results, "APP_NAME")
	if r2.Rotated {
		t.Error("expected APP_NAME to remain unchanged")
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
	}
	opts := Options{Strategy: StrategyBlank, Keys: []string{"APP_NAME"}}
	results := Apply(env, opts)
	r, _ := findResult(results, "APP_NAME")
	if !r.Rotated {
		t.Error("expected APP_NAME to be rotated")
	}
	if r.NewValue != "" {
		t.Errorf("expected blank new value, got %q", r.NewValue)
	}
	r2, _ := findResult(results, "DB_HOST")
	if r2.Rotated {
		t.Error("expected DB_HOST to be unchanged")
	}
}

func TestApply_IncrementStrategy(t *testing.T) {
	env := map[string]string{"SECRET": "original"}
	opts := Options{Strategy: StrategyIncrement, Keys: []string{"SECRET"}}
	results := Apply(env, opts)
	r, _ := findResult(results, "SECRET")
	if r.NewValue != "original_rotated" {
		t.Errorf("unexpected new value: %q", r.NewValue)
	}
}

func TestApply_RedactStrategy(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "hunter2"}
	opts := DefaultOptions()
	results := Apply(env, opts)
	r, _ := findResult(results, "DB_PASSWORD")
	if r.NewValue != "REDACTED_DB_PASSWORD" {
		t.Errorf("unexpected redacted value: %q", r.NewValue)
	}
}

func TestToMap_ReturnsNewValues(t *testing.T) {
	env := map[string]string{"TOKEN": "old"}
	results := Apply(env, DefaultOptions())
	m := ToMap(results)
	if m["TOKEN"] == "old" {
		t.Error("expected rotated value in map")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "s",
		"APP_ENV":    "prod",
	}
	results := Apply(env, DefaultOptions())
	summary := GetSummary(results)
	if summary["total"] != 2 {
		t.Errorf("expected total 2, got %d", summary["total"])
	}
	if summary["rotated"] < 1 {
		t.Error("expected at least 1 rotated")
	}
}
