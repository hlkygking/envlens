package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/wryfi/envlens/internal/transform"
)

func sampleTransformResults() []transform.Result {
	return []transform.Result{
		{Key: "APP_ENV", Original: "production", Value: "PRODUCTION", Applied: []transform.Op{transform.OpUppercase}},
		{Key: "DEBUG", Original: "true", Value: "true", Applied: []transform.Op{}},
	}
}

func TestTransformTextReport_ContainsChanged(t *testing.T) {
	out := TransformTextReport(sampleTransformResults())
	if !strings.Contains(out, "[changed]") {
		t.Error("expected [changed] in output")
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
}

func TestTransformTextReport_Summary(t *testing.T) {
	out := TransformTextReport(sampleTransformResults())
	if !strings.Contains(out, "1 changed") {
		t.Errorf("expected summary with 1 changed, got: %s", out)
	}
}

func TestTransformTextReport_Unchanged(t *testing.T) {
	out := TransformTextReport(sampleTransformResults())
	if !strings.Contains(out, "[unchanged]") {
		t.Error("expected [unchanged] entry")
	}
}

func TestTransformJSONReport_ValidJSON(t *testing.T) {
	out, err := TransformJSONReport(sampleTransformResults())
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := v["summary"]; !ok {
		t.Error("expected summary key in JSON")
	}
	if _, ok := v["results"]; !ok {
		t.Error("expected results key in JSON")
	}
}
