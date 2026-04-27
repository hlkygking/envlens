package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envwatch"
	"github.com/yourorg/envlens/internal/report"
)

func sampleWatchEntries() []envwatch.Entry {
	return []envwatch.Entry{
		{Key: "DB_HOST", OldValue: "", NewValue: "localhost", Status: envwatch.StatusNew, Watched: true},
		{Key: "API_KEY", OldValue: "old", NewValue: "", Status: envwatch.StatusRemoved, Watched: true},
		{Key: "PORT", OldValue: "8080", NewValue: "9090", Status: envwatch.StatusChanged, Watched: true},
		{Key: "DEBUG", OldValue: "true", NewValue: "true", Status: envwatch.StatusUnchanged, Watched: true},
		{Key: "INTERNAL", OldValue: "x", NewValue: "y", Status: envwatch.StatusChanged, Watched: false},
	}
}

func TestEnvWatchTextReport_ContainsSections(t *testing.T) {
	entries := sampleWatchEntries()
	out := report.EnvWatchTextReport(entries, false)

	for _, section := range []string{"New Keys", "Removed Keys", "Changed Keys"} {
		if !strings.Contains(out, section) {
			t.Errorf("expected section %q in output", section)
		}
	}
}

func TestEnvWatchTextReport_HidesUnchanged(t *testing.T) {
	entries := sampleWatchEntries()
	out := report.EnvWatchTextReport(entries, false)
	if strings.Contains(out, "Unchanged Keys") {
		t.Error("unchanged section should be hidden when showUnchanged=false")
	}
}

func TestEnvWatchTextReport_ShowsUnchanged(t *testing.T) {
	entries := sampleWatchEntries()
	out := report.EnvWatchTextReport(entries, true)
	if !strings.Contains(out, "Unchanged Keys") {
		t.Error("expected unchanged section when showUnchanged=true")
	}
}

func TestEnvWatchTextReport_Summary(t *testing.T) {
	entries := sampleWatchEntries()
	out := report.EnvWatchTextReport(entries, false)
	if !strings.Contains(out, "Summary:") {
		t.Error("expected summary line in output")
	}
}

func TestEnvWatchTextReport_NoChanges(t *testing.T) {
	entries := []envwatch.Entry{
		{Key: "A", NewValue: "1", Status: envwatch.StatusUnchanged, Watched: true},
	}
	out := report.EnvWatchTextReport(entries, false)
	if !strings.Contains(out, "No watched changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestEnvWatchJSONReport_ValidJSON(t *testing.T) {
	entries := sampleWatchEntries()
	out, err := report.EnvWatchJSONReport(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected 'entries' key in JSON output")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON output")
	}
}
