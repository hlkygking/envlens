package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/envchain"
)

var sampleChainResults = []envchain.Result{
	{Key: "FOO", Value: "prod_foo", ResolvedBy: "prod", Overridden: true},
	{Key: "BAR", Value: "base_bar", ResolvedBy: "base", Overridden: false},
	{Key: "BAZ", Value: "prod_baz", ResolvedBy: "prod", Overridden: false},
}

func TestEnvChainTextReport_ContainsKeys(t *testing.T) {
	out := EnvChainTextReport(sampleChainResults, envchain.StrategyOverride)
	for _, key := range []string{"FOO", "BAR", "BAZ"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %s in report", key)
		}
	}
}

func TestEnvChainTextReport_ShowsConflict(t *testing.T) {
	out := EnvChainTextReport(sampleChainResults, envchain.StrategyOverride)
	if !strings.Contains(out, "[conflict]") {
		t.Error("expected [conflict] marker in report")
	}
}

func TestEnvChainTextReport_Summary(t *testing.T) {
	out := EnvChainTextReport(sampleChainResults, envchain.StrategyOverride)
	if !strings.Contains(out, "3 keys resolved") {
		t.Error("expected summary line with 3 keys resolved")
	}
	if !strings.Contains(out, "1 conflicts") {
		t.Error("expected 1 conflict in summary")
	}
}

func TestEnvChainTextReport_NoResults(t *testing.T) {
	out := EnvChainTextReport(nil, envchain.StrategyKeep)
	if !strings.Contains(out, "No results") {
		t.Error("expected 'No results' for empty input")
	}
}

func TestEnvChainJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvChainJSONReport(sampleChainResults, envchain.StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["entries"]; !ok {
		t.Error("expected 'entries' field in JSON")
	}
	if parsed["strategy"] != "override" {
		t.Errorf("expected strategy=override, got %v", parsed["strategy"])
	}
}
