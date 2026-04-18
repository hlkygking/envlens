package redact

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

func TestApply_DefaultRules_RedactsSensitive(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret123",
		"APP_NAME":    "myapp",
		"API_KEY":     "abc-xyz",
	}
	results := Apply(env, nil)

	pw := findResult(results, "DB_PASSWORD")
	if pw == nil || !pw.WasRedacted || pw.Redacted != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted")
	}

	app := findResult(results, "APP_NAME")
	if app == nil || app.WasRedacted {
		t.Errorf("expected APP_NAME to not be redacted")
	}

	key := findResult(results, "API_KEY")
	if key == nil || !key.WasRedacted {
		t.Errorf("expected API_KEY to be redacted")
	}
}

func TestApply_CustomRules(t *testing.T) {
	env := map[string]string{
		"MY_CUSTOM_FIELD": "value",
		"OTHER":          "data",
	}
	rules := []Rule{{KeyPattern: "CUSTOM", Replacement: "***"}}
	results := Apply(env, rules)

	custom := findResult(results, "MY_CUSTOM_FIELD")
	if custom == nil || !custom.WasRedacted || custom.Redacted != "***" {
		t.Errorf("expected MY_CUSTOM_FIELD to be redacted with custom rule")
	}

	other := findResult(results, "OTHER")
	if other == nil || other.WasRedacted {
		t.Errorf("expected OTHER to not be redacted")
	}
}

func TestToMap_PreservesRedacted(t *testing.T) {
	env := map[string]string{"TOKEN": "tok123", "HOSTtresults := Apply(env, nil)
	m := ToMap(results)

	if m["TOKEN"] != "[REDACTED]" {
		t.Errorf("expected TOKEN to be redacted in map")
	}
	if m["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged")
	}
}

func TestSummary_Counts(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "s",
		"PLAIN_VAR":  "p",
		"PASSWORD":   "pw",
	}
	results := Apply(env, nil)
	total, redacted := Summary(results)

	if total != 3 {
		t.Errorf("expected total=3, got %d", total)
	}
	if redacted != 2 {
		t.Errorf("expected redacted=2, got %d", redacted)
	}
}
