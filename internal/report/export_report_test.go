package report

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleExportSummary() ExportSummary {
	return ExportSummary{
		Format: "dotenv",
		Path:   "/tmp/out.env",
		Keys:   5,
	}
}

func TestExportTextReport_ContainsFields(t *testing.T) {
	s := sampleExportSummary()
	out := ExportTextReport(s)
	for _, want := range []string{"dotenv", "/tmp/out.env", "5", "Export Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestExportJSONReport_ValidJSON(t *testing.T) {
	s := sampleExportSummary()
	out, err := ExportJSONReport(s)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestExportJSONReport_Fields(t *testing.T) {
	s := sampleExportSummary()
	out, _ := ExportJSONReport(s)
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if m["format"] != "dotenv" {
		t.Errorf("expected format dotenv, got %v", m["format"])
	}
	if int(m["keys"].(float64)) != 5 {
		t.Errorf("expected keys 5, got %v", m["keys"])
	}
}
