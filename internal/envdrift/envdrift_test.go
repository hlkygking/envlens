package envdrift

import (
	"testing"
)

func findResult(results []Result, key string) *Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Match(t *testing.T) {
	base := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	targets := map[string]map[string]string{
		"staging": {"APP_ENV": "production", "PORT": "8080"},
	}
	results := Apply(base, targets)
	for _, r := range results {
		if r.Status != StatusMatch {
			t.Errorf("expected match for %s, got %s", r.Key, r.Status)
		}
	}
}

func TestApply_Drifted(t *testing.T) {
	base := map[string]string{"LOG_LEVEL": "info"}
	targets := map[string]map[string]string{
		"staging": {"LOG_LEVEL": "debug"},
	}
	results := Apply(base, targets)
	r := findResult(results, "LOG_LEVEL")
	if r == nil {
		t.Fatal("expected LOG_LEVEL result")
	}
	if r.Status != StatusDrifted {
		t.Errorf("expected drifted, got %s", r.Status)
	}
	if r.Values["staging"] != "debug" {
		t.Errorf("expected staging=debug, got %s", r.Values["staging"])
	}
}

func TestApply_Missing(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	targets := map[string]map[string]string{
		"prod": {},
	}
	results := Apply(base, targets)
	r := findResult(results, "DB_HOST")
	if r == nil {
		t.Fatal("expected DB_HOST result")
	}
	if r.Status != StatusMissing {
		t.Errorf("expected missing, got %s", r.Status)
	}
}

func TestApply_KeyOnlyInTarget(t *testing.T) {
	base := map[string]string{}
	targets := map[string]map[string]string{
		"prod": {"EXTRA_KEY": "value"},
	}
	results := Apply(base, targets)
	r := findResult(results, "EXTRA_KEY")
	if r == nil {
		t.Fatal("expected EXTRA_KEY result")
	}
	if r.Status != StatusMissing {
		t.Errorf("expected missing for key only in target, got %s", r.Status)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusDrifted},
		{Key: "C", Status: StatusMissing},
		{Key: "D", Status: StatusMatch},
	}
	s := GetSummary(results)
	if s.Total != 4 {
		t.Errorf("expected total=4, got %d", s.Total)
	}
	if s.Match != 2 {
		t.Errorf("expected match=2, got %d", s.Match)
	}
	if s.Drifted != 1 {
		t.Errorf("expected drifted=1, got %d", s.Drifted)
	}
	if s.Missing != 1 {
		t.Errorf("expected missing=1, got %d", s.Missing)
	}
}

func TestApply_MultipleTargets(t *testing.T) {
	base := map[string]string{"REGION": "us-east-1"}
	targets := map[string]map[string]string{
		"eu":   {"REGION": "eu-west-1"},
		"apac": {"REGION": "ap-southeast-1"},
	}
	results := Apply(base, targets)
	r := findResult(results, "REGION")
	if r == nil {
		t.Fatal("expected REGION result")
	}
	if r.Status != StatusDrifted {
		t.Errorf("expected drifted, got %s", r.Status)
	}
	if len(r.Values) != 3 { // baseline + eu + apac
		t.Errorf("expected 3 value entries, got %d", len(r.Values))
	}
}
