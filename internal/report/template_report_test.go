package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/template"
)

func sampleTemplateResults() []template.Result {
	return []template.Result{
		{Key: "URL", Template: "http://{{HOST}}:{{PORT}}", Rendered: "http://localhost:8080", Missing: nil, OK: true},
		{Key: "DSN", Template: "{{USER}}:{{PASS}}@db", Rendered: "<missing:USER>:<missing:PASS>@db", Missing: []string{"USER", "PASS"}, OK: false},
	}
}

func TestTemplateTextReport_ContainsSections(t *testing.T) {
	out := TemplateTextReport(sampleTemplateResults())
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] section")
	}
	if !strings.Contains(out, "[FAILED]") {
		t.Error("expected [FAILED] section")
	}
	if !strings.Contains(out, "URL") {
		t.Error("expected URL in report")
	}
}

func TestTemplateTextReport_Summary(t *testing.T) {
	out := TemplateTextReport(sampleTemplateResults())
	if !strings.Contains(out, "OK: 1") {
		t.Error("expected OK: 1")
	}
	if !strings.Contains(out, "Failed: 1") {
		t.Error("expected Failed: 1")
	}
}

func TestTemplateTextReport_ShowsMissing(t *testing.T) {
	out := TemplateTextReport(sampleTemplateResults())
	if !strings.Contains(out, "USER") {
		t.Error("expected missing key USER in report")
	}
}

func TestTemplateJSONReport_ValidJSON(t *testing.T) {
	out, err := TemplateJSONReport(sampleTemplateResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["results"]; !ok {
		t.Error("expected results key in JSON")
	}
	if _, ok := m["total"]; !ok {
		t.Error("expected total key in JSON")
	}
}
