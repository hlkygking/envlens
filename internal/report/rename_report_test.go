package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/rename"
)

func sampleRenameResults() []rename.Result {
	return []rename.Result{
		{OldKey: "OLD_DB", NewKey: "NEW_DB", Value: "postgres", Renamed: true},
		{OldKey: "MISSING", NewKey: "OTHER", Skipped: true},
		{OldKey: "KEEP", NewKey: "KEEP", Value: "yes"},
	}
}

func TestRenameTextReport_ContainsSections(t *testing.T) {
	out := RenameTextReport(sampleRenameResults())
	for _, want := range []string{"RENAMED", "SKIPPED", "UNCHANGED", "Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestRenameTextReport_Summary(t *testing.T) {
	out := RenameTextReport(sampleRenameResults())
	if !strings.Contains(out, "Renamed:   1") {
		t.Error("expected renamed count 1")
	}
	if !strings.Contains(out, "Skipped:   1") {
		t.Error("expected skipped count 1")
	}
	if !strings.Contains(out, "Unchanged: 1") {
		t.Error("expected unchanged count 1")
	}
}

func TestRenameJSONReport_ValidJSON(t *testing.T) {
	out, err := RenameJSONReport(sampleRenameResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestRenameJSONReport_Fields(t *testing.T) {
	out, _ := RenameJSONReport(sampleRenameResults())
	for _, field := range []string{"results", "renamed", "skipped", "unchanged"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON", field)
		}
	}
}
