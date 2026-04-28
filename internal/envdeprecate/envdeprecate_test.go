package envdeprecate

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

func TestApply_OK(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080"}
	results := Apply(env, nil)
	r := findResult(results, "APP_PORT")
	if r == nil {
		t.Fatal("expected result for APP_PORT")
	}
	if r.Status != StatusOK {
		t.Errorf("expected ok, got %s", r.Status)
	}
}

func TestApply_Deprecated(t *testing.T) {
	env := map[string]string{"OLD_API_KEY": "secret"}
	rules := []Rule{{Pattern: "^OLD_", Reason: "legacy prefix"}}
	results := Apply(env, rules)
	r := findResult(results, "OLD_API_KEY")
	if r == nil {
		t.Fatal("expected result for OLD_API_KEY")
	}
	if r.Status != StatusDeprecated {
		t.Errorf("expected deprecated, got %s", r.Status)
	}
	if r.Reason != "legacy prefix" {
		t.Errorf("unexpected reason: %s", r.Reason)
	}
}

func TestApply_Renamed(t *testing.T) {
	env := map[string]string{"DB_PASS": "hunter2"}
	rules := []Rule{{Pattern: "^DB_PASS$", Replacement: "DB_PASSWORD", Reason: "standardized"}}
	results := Apply(env, rules)
	r := findResult(results, "DB_PASS")
	if r == nil {
		t.Fatal("expected result for DB_PASS")
	}
	if r.Status != StatusRenamed {
		t.Errorf("expected renamed, got %s", r.Status)
	}
	if r.Replacement != "DB_PASSWORD" {
		t.Errorf("unexpected replacement: %s", r.Replacement)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Status: StatusOK},
		{Status: StatusDeprecated},
		{Status: StatusRenamed},
		{Status: StatusDeprecated},
	}
	summary := GetSummary(results)
	if summary["ok"] != 1 {
		t.Errorf("expected 1 ok, got %d", summary["ok"])
	}
	if summary["deprecated"] != 2 {
		t.Errorf("expected 2 deprecated, got %d", summary["deprecated"])
	}
	if summary["renamed"] != 1 {
		t.Errorf("expected 1 renamed, got %d", summary["renamed"])
	}
}

func TestParseRules_Basic(t *testing.T) {
	raw := []string{"^OLD_:NEW_:use new prefix", "LEGACY:"}
	rules, err := ParseRules(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Replacement != "NEW_" {
		t.Errorf("unexpected replacement: %s", rules[0].Replacement)
	}
	if rules[0].Reason != "use new prefix" {
		t.Errorf("unexpected reason: %s", rules[0].Reason)
	}
	if rules[1].Replacement != "" {
		t.Errorf("expected empty replacement for rule 2")
	}
}
