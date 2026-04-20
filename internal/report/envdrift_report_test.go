package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envdrift"
)

func sampleDriftResults() []envdrift.Result {
	return []envdrift.Result{
		{
			Key:      "LOG_LEVEL",
			Status:   envdrift.StatusDrifted,
			Baseline: "info",
			Values:   map[string]string{"baseline": "info", "staging": "debug"},
		},
		{
			Key:      "DB_HOST",
			Status:   envdrift.StatusMissing,
			Baseline: "localhost",
			Values:   map[string]string{"baseline": "localhost", "prod": ""},
		},
		{
			Key:      "PORT",
			Status:   envdrift.StatusMatch,
			Baseline: "8080",
			Values:   map[string]string{"baseline": "8080", "staging": "8080"},
		},
	}
}

func TestEnvDriftTextReport_ContainsSections(t *testing.T) {
	results := sampleDriftResults()
	summary := envdrift.GetSummary(results)
	out := EnvDriftTextReport(results, summary)

	if !strings.Contains(out, "[DRIFTED]") {
		t.Error("expected DRIFTED section")
	}
	if !strings.Contains(out, "[MISSING]") {
		t.Error("expected MISSING section")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
}

func TestEnvDriftTextReport_NoDrift(t *testing.T) {
	results := []envdrift.Result{
		{Key: "PORT", Status: envdrift.StatusMatch, Baseline: "8080", Values: map[string]string{"baseline": "8080"}},
	}
	summary := envdrift.GetSummary(results)
	out := EnvDriftTextReport(results, summary)

	if !strings.Contains(out, "No drift detected") {
		t.Error("expected no-drift message")
	}
}

func TestEnvDriftTextReport_Summary(t *testing.T) {
	results := sampleDriftResults()
	summary := envdrift.GetSummary(results)
	out := EnvDriftTextReport(results, summary)

	if !strings.Contains(out, "total=3") {
		t.Error("expected total=3 in summary")
	}
	if !strings.Contains(out, "drifted=1") {
		t.Error("expected drifted=1 in summary")
	}
	if !strings.Contains(out, "missing=1") {
		t.Error("expected missing=1 in summary")
	}
}

func TestEnvDriftJSONReport_ValidJSON(t *testing.T) {
	results := sampleDriftResults()
	summary := envdrift.GetSummary(results)
	out, err := EnvDriftJSONReport(results, summary)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("expected 'results' key in JSON")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON")
	}
}
