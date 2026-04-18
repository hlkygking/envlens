package rename

import (
	"testing"
)

func findResult(results []Result, oldKey string) *Result {
	for i := range results {
		if results[i].OldKey == oldKey {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Renamed(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value1", "KEEP": "v2"}
	rules := []Rule{{From: "OLD_KEY", To: "NEW_KEY"}}
	results := Apply(env, rules)
	r := findResult(results, "OLD_KEY")
	if r == nil || !r.Renamed || r.NewKey != "NEW_KEY" || r.Value != "value1" {
		t.Errorf("expected renamed result, got %+v", r)
	}
}

func TestApply_Skipped(t *testing.T) {
	env := map[string]string{"KEEP": "v"}
	rules := []Rule{{From: "MISSING", To: "OTHER"}}
	results := Apply(env, rules)
	r := findResult(results, "MISSING")
	if r == nil || !r.Skipped {
		t.Errorf("expected skipped result, got %+v", r)
	}
}

func TestApply_Unchanged(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	rules := []Rule{{From: "A", To: "AA"}}
	results := Apply(env, rules)
	r := findResult(results, "B")
	if r == nil || r.Renamed || r.Skipped || r.NewKey != "B" {
		t.Errorf("expected unchanged B, got %+v", r)
	}
}

func TestToMap(t *testing.T) {
	env := map[string]string{"OLD": "val"}
	rules := []Rule{{From: "OLD", To: "NEW"}}
	m := ToMap(Apply(env, rules))
	if m["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %v", m)
	}
	if _, ok := m["OLD"]; ok {
		t.Error("old key should not exist in ToMap output")
	}
}

func TestSummary(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	rules := []Rule{{From: "A", To: "AA"}, {From: "MISSING", To: "X"}}
	renamed, skipped, unchanged := Summary(Apply(env, rules))
	if renamed != 1 || skipped != 1 || unchanged != 1 {
		t.Errorf("unexpected summary: renamed=%d skipped=%d unchanged=%d", renamed, skipped, unchanged)
	}
}

func TestParseRules(t *testing.T) {
	raw := []string{"OLD=NEW", "FOO = BAR", "invalid"}
	rules := ParseRules(raw)
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[1].From != "FOO" || rules[1].To != "BAR" {
		t.Errorf("unexpected rule: %+v", rules[1])
	}
}
