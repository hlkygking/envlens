package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/envcompact"
	"github.com/yourusername/envlens/internal/report"
)

func sampleCompactResults() []envcompact.Result {
	return []envcompact.Result{
		{Key: "APP_NAME", Value: "myapp", Removed: false, Reason: ""},
		{Key: "EMPTY_KEY", Value: "", Removed: true, Reason: "empty value"},
		{Key: "PORT", Value: "8080", Removed: true, Reason: "matches default"},
	}
}

func TestEnvCompactTextReport_ContainsSections(t *testing.T) {
	out := report.EnvCompactTextReport(sampleCompactResults())
	if !strings.Contains(out, "Removed:") {
		t.Error("expected 'Removed:' section in output")
	}
	if !strings.Contains(out, "EMPTY_KEY") {
		t.Error("expected EMPTY_KEY in output")
	}
	if !strings.Contains(out, "empty value") {
		t.Error("expected reason 'empty value' in output")
	}
}

func TestEnvCompactTextReport_Summary(t *testing.T) {
	out := report.EnvCompactTextReport(sampleCompactResults())
	if !strings.Contains(out, "Kept:    1") {
		t.Errorf("expected kept count 1, got:\n%s", out)
	}
	if !strings.Contains(out, "Removed: 2") {
		t.Errorf("expected removed count 2, got:\n%s", out)
	}
}

func TestEnvCompactTextReport_NoRemoved(t *testing.T) {
	results := []envcompact.Result{
		{Key: "A", Value: "1", Removed: false},
	}
	out := report.EnvCompactTextReport(results)
	if !strings.Contains(out, "No entries removed.") {
		t.Error("expected 'No entries removed.' message")
	}
}

func TestEnvCompactJSONReport_ValidJSON(t *testing.T) {
	out, err := report.EnvCompactJSONReport(sampleCompactResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON output")
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected 'entries' key in JSON output")
	}
}
