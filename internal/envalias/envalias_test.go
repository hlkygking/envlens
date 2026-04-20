package envalias

import (
	"testing"
)

func findResult(results []Result, from string) *Result {
	for i := range results {
		if results[i].OriginalKey == from {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Basic(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	rules := []Rule{{From: "DB_HOST", To: "DATABASE_HOST"}}
	results, out := Apply(env, rules)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := findResult(results, "DB_HOST")
	if r == nil || !r.Applied {
		t.Fatal("expected DB_HOST alias to be applied")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Error("original key DB_HOST should be preserved")
	}
}

func TestApply_SkippedMissingSource(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []Rule{{From: "DB_HOST", To: "DATABASE_HOST"}}
	results, out := Apply(env, rules)

	r := findResult(results, "DB_HOST")
	if r == nil || r.Applied {
		t.Fatal("expected alias to be skipped when source key is missing")
	}
	if _, ok := out["DATABASE_HOST"]; ok {
		t.Error("alias key should not appear when source is missing")
	}
}

func TestApply_Conflict(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DATABASE_HOST": "prod-host"}
	rules := []Rule{{From: "DB_HOST", To: "DATABASE_HOST"}}
	results, out := Apply(env, rules)

	r := findResult(results, "DB_HOST")
	if r == nil || !r.Conflict {
		t.Fatal("expected conflict when alias key already exists")
	}
	if out["DATABASE_HOST"] != "prod-host" {
		t.Error("existing alias key value should not be overwritten on conflict")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Applied: true},
		{Applied: true},
		{Applied: false, Conflict: false}, // skipped
		{Applied: false, Conflict: true},
	}
	s := GetSummary(results)
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Applied != 2 {
		t.Errorf("expected Applied=2, got %d", s.Applied)
	}
	if s.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", s.Skipped)
	}
	if s.Conflicts != 1 {
		t.Errorf("expected Conflicts=1, got %d", s.Conflicts)
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"OLD_KEY=NEW_KEY", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "OLD_KEY" || rules[0].To != "NEW_KEY" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_SkipsMalformed(t *testing.T) {
	rules, _ := ParseRules([]string{"NOEQUALS", "=MISSING_FROM", "MISSING_TO="})
	// "MISSING_TO=" has empty To, so it should be skipped
	// "=MISSING_FROM" has empty From, skipped
	// "NOEQUALS" has no '=', skipped
	if len(rules) != 0 {
		t.Errorf("expected 0 valid rules, got %d", len(rules))
	}
}
