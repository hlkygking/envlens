package validate

import (
	"strings"
	"testing"
)

func TestValidateKeys_Valid(t *testing.T) {
	env := map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "localhost",
		"PORT":       "8080",
	}
	violations := ValidateKeys(env)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestValidateKeys_Invalid(t *testing.T) {
	env := map[string]string{
		"app_env":   "production",
		"1INVALID":  "val",
		"VALID_KEY": "ok",
	}
	violations := ValidateKeys(env)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestValidateRules_Required_Missing(t *testing.T) {
	env := map[string]string{"APP_ENV": "staging"}
	rules := []Rule{
		{Key: "DB_HOST", Required: true},
	}
	violations := ValidateRules(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "missing") {
		t.Errorf("expected missing message, got: %s", violations[0].Message)
	}
}

func TestValidateRules_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	rules := []Rule{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	violations := ValidateRules(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "pattern") {
		t.Errorf("expected pattern message, got: %s", violations[0].Message)
	}
}

func TestValidateRules_PatternPass(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	rules := []Rule{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	violations := ValidateRules(env, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestSummary_NoViolations(t *testing.T) {
	s := Summary(nil)
	if !strings.Contains(s, "passed") {
		t.Errorf("expected passed message, got: %s", s)
	}
}

func TestSummary_WithViolations(t *testing.T) {
	v := []Violation{{Key: "FOO", Message: "something wrong"}}
	s := Summary(v)
	if !strings.Contains(s, "FOO") {
		t.Errorf("expected key in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 violation") {
		t.Errorf("expected count in summary, got: %s", s)
	}
}
