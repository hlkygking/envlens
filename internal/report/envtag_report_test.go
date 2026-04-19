package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/envtag"
)

func sampleTagResults() []envtag.Result {
	return []envtag.Result{
		{Key: "DB_HOST", Value: "localhost", Tags: []string{"database"}, Tagged: true},
		{Key: "AWS_REGION", Value: "us-east-1", Tags: []string{"cloud", "infra"}, Tagged: true},
		{Key: "PORT", Value: "8080", Tags: nil, Tagged: false},
	}
}

func TestEnvTagTextReport_ContainsSections(t *testing.T) {
	out := EnvTagTextReport(sampleTagResults())
	for _, s := range []string{"[Tagged]", "[Untagged]", "DB_HOST", "database", "PORT"} {
		if !strings.Contains(out, s) {
			t.Errorf("expected %q in output", s)
		}
	}
}

func TestEnvTagTextReport_Summary(t *testing.T) {
	out := EnvTagTextReport(sampleTagResults())
	if !strings.Contains(out, "Tagged: 2") {
		t.Error("expected Tagged: 2 in summary")
	}
	if !strings.Contains(out, "Untagged: 1") {
		t.Error("expected Untagged: 1 in summary")
	}
}

func TestEnvTagTextReport_NoResults(t *testing.T) {
	out := EnvTagTextReport([]envtag.Result{})
	if !strings.Contains(out, "(none)") {
		t.Error("expected (none) when no results")
	}
}

func TestEnvTagJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvTagJSONReport(sampleTagResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["summary"]; !ok {
		t.Error("expected summary key in JSON")
	}
	if _, ok := parsed["results"]; !ok {
		t.Error("expected results key in JSON")
	}
}
