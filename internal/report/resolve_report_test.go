package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/envlens/internal/resolve"
)

func sampleResults() []resolve.Result {
	return []resolve.Result{
		{Key: "APP_NAME", Value: "envlens", Source: "file", Resolved: true},
		{Key: "LOG_LEVEL", Value: "info", Source: "default", Resolved: true},
		{Key: "SECRET_KEY", Value: "", Source: "missing", Resolved: false},
	}
}

func TestResolveTextReport_ContainsKeys(t *testing.T) {
	out := ResolveTextReport(sampleResults())
	for _, key := range []string{"APP_NAME", "LOG_LEVEL", "SECRET_KEY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %s in report", key)
		}
	}
}

func TestResolveTextReport_ShowsMissing(t *testing.T) {
	out := ResolveTextReport(sampleResults())
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING status in report")
	}
}

func TestResolveTextReport_ShowsSummary(t *testing.T) {
	out := ResolveTextReport(sampleResults())
	if !strings.Contains(out, "total=") {
		t.Error("expected summary line in report")
	}
}

func TestResolveTextReport_EmptyResults(t *testing.T) {
	out := ResolveTextReport([]resolve.Result{})
	if !strings.Contains(out, "total=0") {
		t.Error("expected total=0 in report for empty results")
	}
}

func TestResolveJSONReport_ValidJSON(t *testing.T) {
	out, err := ResolveJSONReport(sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestResolveJSONReport_Fields(t *testing.T) {
	out, err := ResolveJSONReport(sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, field := range []string{"summary", "results", "key", "source", "resolved"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}
