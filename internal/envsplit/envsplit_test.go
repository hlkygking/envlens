package envsplit_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envsplit"
)

func findResult(results []envsplit.Result, key string) *envsplit.Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_RoutesToBucket(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"APP_PORT": "8080",
		"SECRET":   "abc",
	}
	rules := []envsplit.Rule{
		{Bucket: "database", Pattern: "^DB_"},
		{Bucket: "app", Pattern: "^APP_"},
	}
	results, err := envsplit.Apply(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r := findResult(results, "DB_HOST"); r == nil || r.Bucket != "database" {
		t.Errorf("expected DB_HOST in database bucket")
	}
	if r := findResult(results, "APP_PORT"); r == nil || r.Bucket != "app" {
		t.Errorf("expected APP_PORT in app bucket")
	}
	if r := findResult(results, "SECRET"); r == nil || r.Bucket != "default" {
		t.Errorf("expected SECRET in default bucket")
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	rules := []envsplit.Rule{{Bucket: "bad", Pattern: "[invalid"}}
	_, err := envsplit.Apply(env, rules)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestApply_NoRules_AllDefault(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	results, err := envsplit.Apply(env, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Bucket != "default" {
			t.Errorf("expected default bucket, got %s", r.Bucket)
		}
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []envsplit.Result{
		{Key: "A", Bucket: "alpha"},
		{Key: "B", Bucket: "alpha"},
		{Key: "C", Bucket: "beta"},
	}
	s := envsplit.GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Buckets["alpha"] != 2 {
		t.Errorf("expected alpha=2, got %d", s.Buckets["alpha"])
	}
	if s.Buckets["beta"] != 1 {
		t.Errorf("expected beta=1, got %d", s.Buckets["beta"])
	}
}

func TestFilterByBucket(t *testing.T) {
	results := []envsplit.Result{
		{Key: "A", Bucket: "alpha"},
		{Key: "B", Bucket: "beta"},
		{Key: "C", Bucket: "alpha"},
	}
	filtered := envsplit.FilterByBucket(results, "alpha")
	if len(filtered) != 2 {
		t.Errorf("expected 2 results, got %d", len(filtered))
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	results := []envsplit.Result{
		{Key: "X", Value: "1", Bucket: "a"},
		{Key: "Y", Value: "2", Bucket: "b"},
	}
	m := envsplit.ToMap(results)
	if m["X"] != "1" || m["Y"] != "2" {
		t.Errorf("ToMap produced unexpected values: %v", m)
	}
}
