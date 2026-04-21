package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/envpivot"
	"github.com/yourusername/envlens/internal/report"
)

func samplePivotEntries() []envpivot.Entry {
	return []envpivot.Entry{
		{
			Key:     "PORT",
			Values:  map[string]string{"dev": "3000", "prod": "8080"},
			Missing: nil,
			Uniform: false,
		},
		{
			Key:     "APP",
			Values:  map[string]string{"dev": "myapp", "prod": "myapp"},
			Missing: nil,
			Uniform: true,
		},
		{
			Key:     "SECRET",
			Values:  map[string]string{"dev": "abc"},
			Missing: []string{"prod"},
			Uniform: false,
		},
	}
}

func TestEnvPivotTextReport_ContainsSections(t *testing.T) {
	out := report.EnvPivotTextReport(samplePivotEntries(), []string{"dev", "prod"})
	for _, want := range []string{"PORT", "APP", "SECRET", "DIVERGENT", "uniform"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestEnvPivotTextReport_ShowsMissing(t *testing.T) {
	out := report.EnvPivotTextReport(samplePivotEntries(), []string{"dev", "prod"})
	if !strings.Contains(out, "(missing)") {
		t.Error("expected '(missing)' marker for absent scope value")
	}
}

func TestEnvPivotTextReport_Summary(t *testing.T) {
	out := report.EnvPivotTextReport(samplePivotEntries(), []string{"dev", "prod"})
	if !strings.Contains(out, "total=3") {
		t.Error("expected total=3 in summary")
	}
	if !strings.Contains(out, "uniform=1") {
		t.Error("expected uniform=1 in summary")
	}
}

func TestEnvPivotTextReport_NoEntries(t *testing.T) {
	out := report.EnvPivotTextReport(nil, []string{"dev"})
	if !strings.Contains(out, "No entries") {
		t.Error("expected 'No entries' for empty input")
	}
}

func TestEnvPivotJSONReport_ValidJSON(t *testing.T) {
	out, err := report.EnvPivotJSONReport(samplePivotEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected 'entries' key in JSON")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected 'summary' key in JSON")
	}
}
