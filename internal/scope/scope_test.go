package scope

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

func makeScopes() []Scope {
	return []Scope{
		{Name: "dev", Env: map[string]string{"APP_ENV": "dev", "DB_HOST": "localhost", "SHARED": "same"}},
		{Name: "staging", Env: map[string]string{"APP_ENV": "staging", "DB_HOST": "staging-db", "SHARED": "same"}},
		{Name: "prod", Env: map[string]string{"APP_ENV": "prod", "DB_HOST": "prod-db", "SHARED": "same"}},
	}
}

func TestCompare_DetectsDivergent(t *testing.T) {
	results := Compare(makeScopes())
	r := findResult(results, "APP_ENV")
	if r == nil {
		t.Fatal("expected APP_ENV in results")
	}
	if r.Uniform {
		t.Error("APP_ENV should be divergent across scopes")
	}
}

func TestCompare_DetectsUniform(t *testing.T) {
	results := Compare(makeScopes())
	r := findResult(results, "SHARED")
	if r == nil {
		t.Fatal("expected SHARED in results")
	}
	if !r.Uniform {
		t.Error("SHARED should be uniform across scopes")
	}
}

func TestCompare_AllScopesPresent(t *testing.T) {
	results := Compare(makeScopes())
	r := findResult(results, "DB_HOST")
	if r == nil {
		t.Fatal("expected DB_HOST")
	}
	if len(r.Values) != 3 {
		t.Errorf("expected 3 scope values, got %d", len(r.Values))
	}
}

func TestSummary_Counts(t *testing.T) {
	results := Compare(makeScopes())
	u, d := Summary(results)
	if u+d != len(results) {
		t.Errorf("summary counts don't add up: %d+%d != %d", u, d, len(results))
	}
	if u == 0 {
		t.Error("expected at least one uniform key")
	}
	if d == 0 {
		t.Error("expected at least one divergent key")
	}
}

func TestCompare_EmptyScopes(t *testing.T) {
	results := Compare([]Scope{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
