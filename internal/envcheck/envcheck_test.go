package envcheck

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
	env := map[string]string{"PORT": "8080"}
	rules := []Rule{{Key: "PORT", Required: true}}
	results := Apply(env, rules)
	r := findResult(results, "PORT")
	if r == nil || r.Status != StatusOK {
		t.Fatalf("expected OK for PORT, got %v", r)
	}
}

func TestApply_Missing(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{{Key: "SECRET", Required: true}}
	results := Apply(env, rules)
	r := findResult(results, "SECRET")
	if r == nil || r.Status != StatusMissing {
		t.Fatalf("expected Missing for SECRET, got %v", r)
	}
}

func TestApply_OptionalNotSet(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{{Key: "DEBUG", Required: false}}
	results := Apply(env, rules)
	r := findResult(results, "DEBUG")
	if r == nil || r.Status != StatusOK {
		t.Fatalf("expected OK for optional missing key, got %v", r)
	}
}

func TestApply_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []Rule{{Key: "PORT", Required: true, Pattern: `^\d+$`}}
	results := Apply(env, rules)
	r := findResult(results, "PORT")
	if r == nil || r.Status != StatusOK {
		t.Fatalf("expected OK for PORT with matching pattern, got %v", r)
	}
}

func TestApply_PatternFail(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	rules := []Rule{{Key: "PORT", Required: true, Pattern: `^\d+$`}}
	results := Apply(env, rules)
	r := findResult(results, "PORT")
	if r == nil || r.Status != StatusInvalid {
		t.Fatalf("expected Invalid for PORT with bad value, got %v", r)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Status: StatusOK},
		{Key: "B", Status: StatusMissing},
		{Key: "C", Status: StatusInvalid},
		{Key: "D", Status: StatusOK},
	}
	s := GetSummary(results)
	if s.Total != 4 || s.OK != 2 || s.Missing != 1 || s.Invalid != 1 {
		t.Fatalf("unexpected summary: %+v", s)
	}
}

func TestParseRules_Basic(t *testing.T) {
	lines := []string{
		"PORT:required:^\\d+$",
		"DEBUG:optional",
		"SECRET:required",
		"# comment",
		"",
	}
	rules := ParseRules(lines)
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if !rules[0].Required || rules[0].Pattern == "" {
		t.Errorf("PORT rule malformed: %+v", rules[0])
	}
	if rules[1].Required {
		t.Errorf("DEBUG should not be required")
	}
}
