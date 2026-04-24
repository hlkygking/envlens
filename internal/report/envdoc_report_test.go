package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envdoc"
)

var sampleDocResults = []envdoc.Result{
	{Key: "APP_PORT", Found: true, Description: "HTTP port", Required: true, Default: "3000"},
	{Key: "DB_URL", Found: true, Description: "Database URL", Required: true},
	{Key: "ORPHAN_KEY", Found: false},
	{Key: "ANOTHER_ORPHAN", Found: false},
}

func TestEnvDocTextReport_ContainsSections(t *testing.T) {
	out := EnvDocTextReport(sampleDocResults, true)
	if !strings.Contains(out, "=== Documented ===") {
		t.Error("expected Documented section")
	}
	if !strings.Contains(out, "=== Undocumented ===") {
		t.Error("expected Undocumented section")
	}
}

func TestEnvDocTextReport_ShowsDescription(t *testing.T) {
	out := EnvDocTextReport(sampleDocResults, true)
	if !strings.Contains(out, "HTTP port") {
		t.Error("expected description 'HTTP port' in output")
	}
	if !strings.Contains(out, "[required]") {
		t.Error("expected [required] marker")
	}
	if !strings.Contains(out, "default: 3000") {
		t.Error("expected default value in output")
	}
}

func TestEnvDocTextReport_ShowsUndocumented(t *testing.T) {
	out := EnvDocTextReport(sampleDocResults, true)
	if !strings.Contains(out, "ORPHAN_KEY") {
		t.Error("expected ORPHAN_KEY in undocumented section")
	}
}

func TestEnvDocTextReport_Summary(t *testing.T) {
	out := EnvDocTextReport(sampleDocResults, true)
	if !strings.Contains(out, "total=4") {
		t.Error("expected total=4 in summary")
	}
	if !strings.Contains(out, "documented=2") {
		t.Error("expected documented=2 in summary")
	}
	if !strings.Contains(out, "undocumented=2") {
		t.Error("expected undocumented=2 in summary")
	}
}

func TestEnvDocJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvDocJSONReport(sampleDocResults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected 'entries' key in JSON")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON")
	}
}
