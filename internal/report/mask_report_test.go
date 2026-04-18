package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "envlens",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "key123",
		"PORT":        "8080",
	}
}

func TestMaskTextReport_ContainsMasked(t *testing.T) {
	var buf bytes.Buffer
	MaskTextReport(sampleEnv(), &buf)
	out := buf.String()

	if strings.Contains(out, "supersecret") {
		t.Error("DB_PASSWORD plain value should not appear in output")
	}
	if strings.Contains(out, "key123") {
		t.Error("API_KEY plain value should not appear in output")
	}
	if !strings.Contains(out, "envlens") {
		t.Error("APP_NAME value should appear unmasked")
	}
	if !strings.Contains(out, "[masked]") {
		t.Error("output should contain [masked] marker")
	}
}

func TestMaskTextReport_Summary(t *testing.T) {
	var buf bytes.Buffer
	MaskTextReport(sampleEnv(), &buf)
	if !strings.Contains(buf.String(), "Total: 4") {
		t.Error("summary line should show total key count")
	}
}

func TestMaskJSONReport_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := MaskJSONReport(sampleEnv(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entries []MaskedEntry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}
}

func TestMaskJSONReport_MaskedFlag(t *testing.T) {
	var buf bytes.Buffer
	MaskJSONReport(sampleEnv(), &buf)
	var entries []MaskedEntry
	json.Unmarshal(buf.Bytes(), &entries)

	for _, e := range entries {
		if e.Key == "DB_PASSWORD" && !e.Masked {
			t.Error("DB_PASSWORD should have Masked=true")
		}
		if e.Key == "APP_NAME" && e.Masked {
			t.Error("APP_NAME should have Masked=false")
		}
	}
}
