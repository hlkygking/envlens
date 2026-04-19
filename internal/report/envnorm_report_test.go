package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envnorm"
)

func sampleNormResults() []envnorm.Result {
	return []envnorm.Result{
		{OriginalKey: "db_host", NormalizedKey: "DB_HOST", Value: "localhost", Changed: true},
		{OriginalKey: "APP_PORT", NormalizedKey: "APP_PORT", Value: "8080", Changed: false},
		{OriginalKey: "api_key", NormalizedKey: "API_KEY", Value: "secret", Changed: true},
	}
}

func TestEnvNormTextReport_ContainsSections(t *testing.T) {
	out := EnvNormTextReport(sampleNormResults())
	if !strings.Contains(out, "CHANGED") {
		t.Error("expected CHANGED section")
	}
	if !strings.Contains(out, "UNCHANGED") {
		t.Error("expected UNCHANGED section")
	}
}

func TestEnvNormTextReport_ShowsRename(t *testing.T) {
	out := EnvNormTextReport(sampleNormResults())
	if !strings.Contains(out, "db_host -> DB_HOST") {
		t.Error("expected rename entry db_host -> DB_HOST")
	}
}

func TestEnvNormTextReport_Summary(t *testing.T) {
	out := EnvNormTextReport(sampleNormResults())
	if !strings.Contains(out, "3 total") {
		t.Error("expected '3 total' in summary")
	}
	if !strings.Contains(out, "2 changed") {
		t.Error("expected '2 changed' in summary")
	}
}

func TestEnvNormJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvNormJSONReport(sampleNormResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["results"]; !ok {
		t.Error("expected 'results' key")
	}
	if _, ok := m["summary"]; !ok {
		t.Error("expected 'summary' key")
	}
}
