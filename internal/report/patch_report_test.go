package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/patch"
)

func samplePatchResults() []patch.Result {
	return []patch.Result{
		{Rule: patch.Rule{Op: patch.OpSet, Key: "FOO", Value: "bar"}, Applied: true, Note: "set"},
		{Rule: patch.Rule{Op: patch.OpDelete, Key: "OLD"}, Applied: true, Note: "deleted"},
		{Rule: patch.Rule{Op: patch.OpRename, Key: "X", To: "Y"}, Applied: true, Note: "X -> Y"},
		{Rule: patch.Rule{Op: patch.OpDelete, Key: "MISSING"}, Applied: false, Note: "key not found"},
	}
}

func TestPatchTextReport_ContainsSections(t *testing.T) {
	out := PatchTextReport(samplePatchResults())
	if !strings.Contains(out, "Patch Results") {
		t.Error("expected header")
	}
	if !strings.Contains(out, "FOO") {
		t.Error("expected FOO in output")
	}
	if !strings.Contains(out, "SKIP") {
		t.Error("expected SKIP status")
	}
}

func TestPatchTextReport_Summary(t *testing.T) {
	out := PatchTextReport(samplePatchResults())
	if !strings.Contains(out, "3 applied") {
		t.Errorf("expected 3 applied in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected 1 skipped in summary, got: %s", out)
	}
}

func TestPatchTextReport_NoResults(t *testing.T) {
	out := PatchTextReport([]patch.Result{})
	if !strings.Contains(out, "no rules applied") {
		t.Error("expected empty message")
	}
}

func TestPatchJSONReport_ValidJSON(t *testing.T) {
	out, err := PatchJSONReport(samplePatchResults())
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["results"]; !ok {
		t.Error("expected results key")
	}
	if m["applied"].(float64) != 3 {
		t.Errorf("expected applied=3")
	}
}
