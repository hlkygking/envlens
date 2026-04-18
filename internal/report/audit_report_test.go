package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/envlens/envlens/internal/audit"
)

func sampleFindings() []audit.Finding {
	return []audit.Finding{
		{Key: "DB_PASSWORD", Severity: audit.SeverityHigh, Message: "Sensitive key added in target environment"},
		{Key: "AUTH_TOKEN", Severity: audit.SeverityMedium, Message: "Sensitive key value changed between environments"},
		{Key: "LOG_LEVEL", Severity: audit.SeverityLow, Message: "Key transitioned from empty value to non-empty"},
	}
}

func TestAuditTextReport_ContainsSeverities(t *testing.T) {
	var buf bytes.Buffer
	AuditTextReport(&buf, sampleFindings())
	out := buf.String()
	for _, want := range []string{"HIGH", "MEDIUM", "LOW", "DB_PASSWORD", "AUTH_TOKEN"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestAuditTextReport_NoFindings(t *testing.T) {
	var buf bytes.Buffer
	AuditTextReport(&buf, nil)
	if !strings.Contains(buf.String(), "No audit findings") {
		t.Error("expected no-findings message")
	}
}

func TestAuditTextReport_Summary(t *testing.T) {
	var buf bytes.Buffer
	AuditTextReport(&buf, sampleFindings())
	if !strings.Contains(buf.String(), "3 finding(s)") {
		t.Error("expected summary count of 3")
	}
}

func TestAuditJSONReport_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := AuditJSONReport(&buf, sampleFindings()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
	if out[0]["severity"] != "HIGH" {
		t.Errorf("expected first entry severity HIGH, got %s", out[0]["severity"])
	}
}
