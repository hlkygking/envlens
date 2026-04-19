package envwhere_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/envwhere"
)

func findResult(results []envwhere.Result, key string) *envwhere.Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_SECRET":  "s3cr3t",
	"APP_VERSION": "1.0.0",
	"LOG_LEVEL":   "info",
}

func TestApply_Exact(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "DB_HOST", envwhere.Options{Mode: envwhere.MatchExact})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value != "localhost" {
		t.Errorf("expected localhost, got %s", results[0].Value)
	}
}

func TestApply_Prefix(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "DB_", envwhere.Options{Mode: envwhere.MatchPrefix})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Matched {
			t.Errorf("expected all matched")
		}
	}
}

func TestApply_Suffix(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "_HOST", envwhere.Options{Mode: envwhere.MatchSuffix})
	r := findResult(results, "DB_HOST")
	if r == nil || !r.Matched {
		t.Error("expected DB_HOST to match suffix _HOST")
	}
}

func TestApply_Regex(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "^APP_", envwhere.Options{Mode: envwhere.MatchRegex})
	if len(results) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(results))
	}
}

func TestApply_InvalidRegex(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "[invalid", envwhere.Options{Mode: envwhere.MatchRegex})
	if len(results) != 1 || results[0].Error == "" {
		t.Error("expected error result for invalid regex")
	}
}

func TestApply_CaseFold(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "db_", envwhere.Options{Mode: envwhere.MatchPrefix, CaseFold: true})
	if len(results) != 2 {
		t.Errorf("expected 2 case-folded results, got %d", len(results))
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := envwhere.Apply(sampleEnv, "APP_", envwhere.Options{Mode: envwhere.MatchPrefix})
	matched, total := envwhere.GetSummary(results)
	if matched != total || matched == 0 {
		t.Errorf("unexpected summary: matched=%d total=%d", matched, total)
	}
}
