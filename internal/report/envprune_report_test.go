package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envprune"
	"github.com/yourorg/envlens/internal/report"
)

func samplePruneResults() []envprune.Result {
	return []envprune.Result{
		{Key: "DEBUG_MODE", Value: "true", Pruned: true, Reason: "matches prefix: DEBUG_"},
		{Key: "APP_NAME", Value: "envlens", Pruned: false},
		{Key: "EMPTY_KEY", Value: "", Pruned: true, Reason: "empty value"},
		{Key: "APP_PORT", Value: "8080", Pruned: false},
	}
}

func TestEnvPruneTextReport_ContainsSections(t *testing.T) {
	out := report.EnvPruneTextReport(samplePruneResults())
	for _, section := range []string{"Pruned Keys", "Retained Keys", "Summary"} {
		if !strings.Contains(out, section) {
			t.Errorf("expected section %q in output", section)
		}
	}
}

func TestEnvPruneTextReport_ShowsPrunedKey(t *testing.T) {
	out := report.EnvPruneTextReport(samplePruneResults())
	if !strings.Contains(out, "DEBUG_MODE") {
		t.Error("expected DEBUG_MODE in pruned section")
	}
	if !strings.Contains(out, "empty value") {
		t.Error("expected reason 'empty value' in output")
	}
}

func TestEnvPruneTextReport_Summary(t *testing.T) {
	out := report.EnvPruneTextReport(samplePruneResults())
	if !strings.Contains(out, "2 pruned") {
		t.Error("expected '2 pruned' in summary")
	}
	if !strings.Contains(out, "2 retained") {
		t.Error("expected '2 retained' in summary")
	}
}

func TestEnvPruneTextReport_NoPruned(t *testing.T) {
	results := []envprune.Result{
		{Key: "A", Value: "1", Pruned: false},
	}
	out := report.EnvPruneTextReport(results)
	if !strings.Contains(out, "(none)") {
		t.Error("expected '(none)' when no keys are pruned")
	}
}

func TestEnvPruneJSONReport_ValidJSON(t *testing.T) {
	out, err := report.EnvPruneJSONReport(samplePruneResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, key := range []string{"results", "pruned", "retained"} {
		if _, ok := payload[key]; !ok {
			t.Errorf("expected key %q in JSON output", key)
		}
	}
}
