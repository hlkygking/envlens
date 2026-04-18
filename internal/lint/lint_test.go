package lint

import (
	"testing"
)

func findIssue(issues []Issue, key string) *Issue {
	for i := range issues {
		if issues[i].Key == key {
			return &issues[i]
		}
	}
	return nil
}

func TestLint_LowercaseKey(t *testing.T) {
	issues := Lint(map[string]string{"myKey": "value"})
	if findIssue(issues, "myKey") == nil {
		t.Error("expected issue for lowercase key")
	}
}

func TestLint_InvalidKeyChars(t *testing.T) {
	issues := Lint(map[string]string{"MY-KEY": "value"})
	iss := findIssue(issues, "MY-KEY")
	if iss == nil {
		t.Fatal("expected issue for invalid key char")
	}
	if iss.Severity != "error" {
		t.Errorf("expected error severity, got %s", iss.Severity)
	}
}

func TestLint_LeadingWhitespace(t *testing.T) {
	issues := Lint(map[string]string{"MY_KEY": " value"})
	iss := findIssue(issues, "MY_KEY")
	if iss == nil {
		t.Error("expected issue for leading whitespace")
	}
}

func TestLint_DoubleSpace(t *testing.T) {
	issues := Lint(map[string]string{"MY_KEY": "hello  world"})
	iss := findIssue(issues, "MY_KEY")
	if iss == nil {
		t.Error("expected issue for double space in value")
	}
}

func TestLint_CleanEnv(t *testing.T) {
	issues := Lint(map[string]string{
		"APP_ENV":  "production",
		"LOG_LEVEL": "info",
	})
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestSummary(t *testing.T) {
	issues := []Issue{
		{Key: "A", Severity: "warn"},
		{Key: "B", Severity: "error"},
		{Key: "C", Severity: "warn"},
	}
	w, e := Summary(issues)
	if w != 2 || e != 1 {
		t.Errorf("expected 2 warns 1 error, got %d/%d", w, e)
	}
}
