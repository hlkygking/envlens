package envduplicate

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

func TestApply_NoDuplicates(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"ENV":  "production",
	}
	results := Apply(env)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != StatusUnique {
			t.Errorf("key %s should be unique", r.Key)
		}
	}
}

func TestApply_DetectsDuplicates(t *testing.T) {
	env := map[string]string{
		"DB_PASS":  "secret",
		"API_KEY":  "secret",
		"APP_NAME": "myapp",
	}
	results := Apply(env)

	dbPass, ok := findResult(results, "DB_PASS")
	if !ok {
		t.Fatal("DB_PASS not found")
	}
	if dbPass.Status != StatusDuplicate {
		t.Errorf("DB_PASS should be duplicate")
	}
	if len(dbPass.SharedWith) != 1 || dbPass.SharedWith[0] != "API_KEY" {
		t.Errorf("DB_PASS should share with API_KEY, got %v", dbPass.SharedWith)
	}

	appName, ok := findResult(results, "APP_NAME")
	if !ok {
		t.Fatal("APP_NAME not found")
	}
	if appName.Status != StatusUnique {
		t.Errorf("APP_NAME should be unique")
	}
}

func TestApply_MultipleSharing(t *testing.T) {
	env := map[string]string{
		"A": "same",
		"B": "same",
		"C": "same",
	}
	results := Apply(env)
	r, ok := findResult(results, "A")
	if !ok {
		t.Fatal("A not found")
	}
	if r.Status != StatusDuplicate {
		t.Errorf("A should be duplicate")
	}
	if len(r.SharedWith) != 2 {
		t.Errorf("A should share with 2 keys, got %d", len(r.SharedWith))
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"X": "val1",
		"Y": "val1",
		"Z": "val2",
	}
	results := Apply(env)
	s := GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Duplicates != 2 {
		t.Errorf("expected 2 duplicates, got %d", s.Duplicates)
	}
	if s.Unique != 1 {
		t.Errorf("expected 1 unique, got %d", s.Unique)
	}
	if s.Groups != 1 {
		t.Errorf("expected 1 group, got %d", s.Groups)
	}
}

func TestApply_EmptyMap(t *testing.T) {
	results := Apply(map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results for empty input")
	}
}
