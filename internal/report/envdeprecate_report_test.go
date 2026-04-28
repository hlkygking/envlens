package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envdeprecate"
)

func sampleDeprecateResults() []envdeprecate.Result {
	return []envdeprecate.Result{
		{Key: "OLD_TOKEN", Value: "abc", Status: envdeprecate.StatusDeprecated, Reason: "legacy"},
		{Key: "DB_PASS", Value: "x", Status: envdeprecate.StatusRenamed, Replacement: "DB_PASSWORD", Reason: "standardized"},
		{Key: "APP_PORT", Value: "8080", Status: envdeprecate.StatusOK},
	}
}

func TestEnvDeprecateTextReport_ContainsSections(t *testing.T) {
	out := EnvDeprecateTextReport(sampleDeprecateResults(), true)
	for _, want := range []string{"Deprecated", "Renamed", "OK", "OLD_TOKEN", "DB_PASS", "DB_PASSWORD"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestEnvDeprecateTextReport_ShowsReason(t *testing.T) {
	out := EnvDeprecateTextReport(sampleDeprecateResults(), false)
	if !strings.Contains(out, "legacy") {
		t.Error("expected reason 'legacy' in output")
	}
	if !strings.Contains(out, "standardized") {
		t.Error("expected reason 'standardized' in output")
	}
}

func TestEnvDeprecateTextReport_HidesOK(t *testing.T) {
	out := EnvDeprecateTextReport(sampleDeprecateResults(), false)
	if strings.Contains(out, "=== OK ===") {
		t.Error("expected OK section to be hidden when showOK=false")
	}
}

func TestEnvDeprecateTextReport_Summary(t *testing.T) {
	out := EnvDeprecateTextReport(sampleDeprecateResults(), false)
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line in output")
	}
	if !strings.Contains(out, "1 deprecated") {
		t.Error("expected '1 deprecated' in summary")
	}
}

func TestEnvDeprecateJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvDeprecateJSONReport(sampleDeprecateResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := payload["results"]; !ok {
		t.Error("expected 'results' key in JSON")
	}
	if _, ok := payload["summary"]; !ok {
		t.Error("expected 'summary' key in JSON")
	}
}
