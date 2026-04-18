package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/interpolate"
)

func sampleInterpolateResults() []interpolate.Result {
	return []interpolate.Result{
		{Key: "BASE", Original: "/app", Resolved: "/app", Refs: nil, Missing: nil},
		{Key: "DATA", Original: "${BASE}/data", Resolved: "/app/data", Refs: []string{"BASE"}, Missing: nil},
		{Key: "BAD", Original: "${NOPE}/x", Resolved: "${NOPE}/x", Refs: []string{"NOPE"}, Missing: []string{"NOPE"}},
	}
}

func TestInterpolateTextReport_ContainsSections(t *testing.T) {
	results := sampleInterpolateResults()
	out := InterpolateTextReport(results, nil)
	if !strings.Contains(out, "Interpolation Report") {
		t.Error("expected header")
	}
	if !strings.Contains(out, "DATA") {
		t.Error("expected DATA key")
	}
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING status")
	}
}

func TestInterpolateTextReport_ShowsCycles(t *testing.T) {
	results := sampleInterpolateResults()
	out := InterpolateTextReport(results, []string{"A <-> B"})
	if !strings.Contains(out, "CYCLES") {
		t.Error("expected cycle section")
	}
	if !strings.Contains(out, "A <-> B") {
		t.Error("expected cycle detail")
	}
}

func TestInterpolateTextReport_NoRefs(t *testing.T) {
	results := []interpolate.Result{
		{Key: "PLAIN", Original: "hello", Resolved: "hello"},
	}
	out := InterpolateTextReport(results, nil)
	if !strings.Contains(out, "No interpolation") {
		t.Error("expected no-refs message")
	}
}

func TestInterpolateJSONReport_ValidJSON(t *testing.T) {
	results := sampleInterpolateResults()
	out, err := InterpolateJSONReport(results, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected summary key")
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("expected results key")
	}
}
