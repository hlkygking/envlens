package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/envlens/internal/promote"
)

func samplePromoteResults() []promote.Result {
	return []promote.Result{
		{Key: "NEW_KEY", Value: "val1", Status: "promoted"},
		{Key: "CONFLICT_KEY", Value: "new", Status: "promoted", Conflict: true, Message: "overrode \"old\""},
		{Key: "SAME_KEY", Value: "same", Status: "skipped", Message: "identical"},
	}
}

func TestPromoteTextReport_ContainsSections(t *testing.T) {
	results := samplePromoteResults()
	summary := promote.GetSummary(results)
	out := PromoteTextReport(results, summary)

	if !strings.Contains(out, "PROMOTED:") {
		t.Error("expected PROMOTED section")
	}
	if !strings.Contains(out, "SKIPPED:") {
		t.Error("expected SKIPPED section")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
}

func TestPromoteTextReport_ShowsConflictMessage(t *testing.T) {
	results := samplePromoteResults()
	summary := promote.GetSummary(results)
	out := PromoteTextReport(results, summary)
	if !strings.Contains(out, "overrode") {
		t.Error("expected conflict message in output")
	}
}

func TestPromoteTextReport_Summary(t *testing.T) {
	results := samplePromoteResults()
	summary := promote.GetSummary(results)
	out := PromoteTextReport(results, summary)
	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line")
	}
}

func TestPromoteTextReport_NoChanges(t *testing.T) {
	out := PromoteTextReport([]promote.Result{}, promote.Summary{})
	if !strings.Contains(out, "No changes.") {
		t.Error("expected No changes message")
	}
}

func TestPromoteJSONReport_ValidJSON(t *testing.T) {
	results := samplePromoteResults()
	summary := promote.GetSummary(results)
	out, err := PromoteJSONReport(results, summary)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("expected results key")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected summary key")
	}
}
