package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envreplace"
)

func sampleReplaceResults() []envreplace.Result {
	return []envreplace.Result{
		{Key: "ENV", OldValue: "staging", NewValue: "production", Status: envreplace.StatusReplaced, MatchedRule: "staging"},
		{Key: "HOST", OldValue: "localhost", NewValue: "localhost", Status: envreplace.StatusUnchanged},
		{Key: "BAD", OldValue: "[x", NewValue: "[x", Status: envreplace.StatusError, Error: "invalid regex", MatchedRule: "[x"},
	}
}

func TestEnvReplaceTextReport_ContainsSections(t *testing.T) {
	out := EnvReplaceTextReport(sampleReplaceResults(), false)
	if !strings.Contains(out, "=== Replaced ===") {
		t.Error("expected Replaced section")
	}
	if !strings.Contains(out, "=== Errors ===") {
		t.Error("expected Errors section")
	}
	if strings.Contains(out, "=== Unchanged ===") {
		t.Error("expected Unchanged section to be hidden")
	}
}

func TestEnvReplaceTextReport_ShowsUnchanged(t *testing.T) {
	out := EnvReplaceTextReport(sampleReplaceResults(), true)
	if !strings.Contains(out, "=== Unchanged ===") {
		t.Error("expected Unchanged section when showUnchanged=true")
	}
}

func TestEnvReplaceTextReport_Summary(t *testing.T) {
	out := EnvReplaceTextReport(sampleReplaceResults(), false)
	if !strings.Contains(out, "1 replaced") {
		t.Error("expected 1 replaced in summary")
	}
	if !strings.Contains(out, "1 unchanged") {
		t.Error("expected 1 unchanged in summary")
	}
	if !strings.Contains(out, "1 errors") {
		t.Error("expected 1 errors in summary")
	}
}

func TestEnvReplaceJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvReplaceJSONReport(sampleReplaceResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected entries key in JSON")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected summary key in JSON")
	}
}

func TestEnvReplaceJSONReport_Fields(t *testing.T) {
	out, err := EnvReplaceJSONReport(sampleReplaceResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"status\"") {
		t.Error("expected status field in JSON")
	}
	if !strings.Contains(out, "\"old_value\"") {
		t.Error("expected old_value field in JSON")
	}
	if !strings.Contains(out, "\"matched_rule\"") {
		t.Error("expected matched_rule field in JSON")
	}
}
