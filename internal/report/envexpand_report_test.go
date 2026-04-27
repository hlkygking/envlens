package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envexpand"
	"github.com/yourorg/envlens/internal/report"
)

var sampleExpandResults = []envexpand.Result{
	{Key: "BASE", Original: "/app", Expanded: "/app", Refs: nil, Status: "unchanged"},
	{Key: "LOG_DIR", Original: "${BASE}/logs", Expanded: "/app/logs", Refs: []string{"BASE"}, Status: "ok"},
	{Key: "DB_URL", Original: "postgres://${HOST}/db", Expanded: "postgres://${HOST}/db", Refs: []string{"HOST"}, Status: "unresolved"},
}

func TestEnvExpandTextReport_ContainsSections(t *testing.T) {
	out := report.EnvExpandTextReport(sampleExpandResults, false)
	if !strings.Contains(out, "Expanded") {
		t.Error("expected Expanded section")
	}
	if !strings.Contains(out, "Unresolved") {
		t.Error("expected Unresolved section")
	}
}

func TestEnvExpandTextReport_ShowsUnchanged(t *testing.T) {
	out := report.EnvExpandTextReport(sampleExpandResults, true)
	if !strings.Contains(out, "Unchanged") {
		t.Error("expected Unchanged section when showUnchanged=true")
	}
}

func TestEnvExpandTextReport_HidesUnchanged(t *testing.T) {
	out := report.EnvExpandTextReport(sampleExpandResults, false)
	if strings.Contains(out, "=== Unchanged ===") {
		t.Error("expected Unchanged section to be hidden when showUnchanged=false")
	}
}

func TestEnvExpandTextReport_Summary(t *testing.T) {
	out := report.EnvExpandTextReport(sampleExpandResults, false)
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line")
	}
	if !strings.Contains(out, "total=3") {
		t.Error("expected total=3 in summary")
	}
}

func TestEnvExpandJSONReport_ValidJSON(t *testing.T) {
	out, err := report.EnvExpandJSONReport(sampleExpandResults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("expected results key in JSON")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected summary key in JSON")
	}
}
