package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/envgroup"
)

func sampleGroupResult() envgroup.Result {
	return envgroup.Result{
		Groups: []envgroup.Group{
			{Name: "database", Pattern: "^DB_", Keys: []string{"DB_HOST", "DB_PORT"}},
			{Name: "app", Pattern: "^APP_", Keys: []string{"APP_NAME"}},
		},
		Ungrouped: []string{"LOG_LEVEL"},
	}
}

func TestEnvGroupTextReport_ContainsSections(t *testing.T) {
	out := EnvGroupTextReport(sampleGroupResult())
	for _, want := range []string{"database", "app", "ungrouped", "DB_HOST", "APP_NAME", "LOG_LEVEL"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestEnvGroupTextReport_Summary(t *testing.T) {
	out := EnvGroupTextReport(sampleGroupResult())
	if !strings.Contains(out, "Summary") {
		t.Error("expected Summary section")
	}
}

func TestEnvGroupJSONReport_ValidJSON(t *testing.T) {
	out, err := EnvGroupJSONReport(sampleGroupResult())
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, key := range []string{"groups", "ungrouped", "summary"} {
		if _, ok := v[key]; !ok {
			t.Errorf("missing key %q in JSON", key)
		}
	}
}

func TestEnvGroupJSONReport_GroupCount(t *testing.T) {
	out, err := EnvGroupJSONReport(sampleGroupResult())
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	json.Unmarshal([]byte(out), &v)
	groups := v["groups"].([]interface{})
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}
