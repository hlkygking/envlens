package filter

import (
	"testing"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_PASSWORD":  "secret",
	"DB_HOST":      "db.local",
	"LOG_LEVEL":    "info",
	"FEATURE_FLAG": "true",
}

func findResult(results []Result, key string) *Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Prefix(t *testing.T) {
	res, err := Apply(sampleEnv, Options{Prefix: "APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if findResult(res, "APP_HOST") == nil || findResult(res, "APP_PORT") == nil {
		t.Error("expected APP_HOST and APP_PORT in results")
	}
}

func TestApply_Suffix(t *testing.T) {
	res, err := Apply(sampleEnv, Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestApply_Pattern(t *testing.T) {
	res, err := Apply(sampleEnv, Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := Apply(sampleEnv, Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestApply_Keys(t *testing.T) {
	res, err := Apply(sampleEnv, Options{Keys: []string{"LOG_LEVEL", "FEATURE_FLAG"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
}

func TestApply_NoOptions(t *testing.T) {
	res, err := Apply(sampleEnv, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != len(sampleEnv) {
		t.Fatalf("expected all %d entries, got %d", len(sampleEnv), len(res))
	}
}

func TestGetSummary(t *testing.T) {
	res, _ := Apply(sampleEnv, Options{Prefix: "APP_"})
	s := GetSummary(sampleEnv, res)
	if s.Total != 6 || s.Matched != 2 {
		t.Errorf("unexpected summary: %+v", s)
	}
}
