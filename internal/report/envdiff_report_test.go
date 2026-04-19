package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/envdiff"
)

func sampleDiffEntries() []envdiff.Entry {
	return []envdiff.Entry{
		{Key: "APP_ENV", BaseVal: "dev", TargetVal: "prod", Change: envdiff.Modified},
		{Key: "NEW_KEY", TargetVal: "hello", Change: envdiff.Added},
		{Key: "OLD_KEY", BaseVal: "bye", Change: envdiff.Removed},
		{Key: "STABLE", BaseVal: "same", TargetVal: "same", Change: envdiff.Unchanged},
	}
}

func TestEnvDiffTextReport_ContainsChanges(t *testing.T) {
	out := EnvDiffTextReport(sampleDiffEntries())
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Error("expected added key in output")
	}
	if !strings.Contains(out, "- OLD_KEY") {
		t.Error("expected removed key in output")
	}
	if !strings.Contains(out, "~ APP_ENV") {
		t.Error("expected modified key in output")
	}
}

func TestEnvDiffTextReport_HidesUnchanged(t *testing.T) {
	out := EnvDiffTextReport(sampleDiffEntries())
	if strings.Contains(out, "STABLE") {
		t.Error("unchanged keys should not appear in text report")
	}
}

func TestEnvDiffTextReport_Summary(t *testing.T) {
	out := EnvDiffTextReport(sampleDiffEntries())
	if !strings.Contains(out, "Summary:") {
		t.Error("expected summary line")
	}
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected added count, got: %s", out)
	}
}

func TestEnvDiffJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvDiffJSONReport(sampleDiffEntries())
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
