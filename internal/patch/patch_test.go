package patch

import (
	"testing"
)

func findResult(results []Result, key string) *Result {
	for i := range results {
		if results[i].Rule.Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Set(t *testing.T) {
	env := map[string]string{"A": "1"}
	out, results := Apply(env, []Rule{{Op: OpSet, Key: "B", Value: "2"}})
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
	if r := findResult(results, "B"); r == nil || !r.Applied {
		t.Error("expected applied result for B")
	}
}

func TestApply_Delete(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	out, results := Apply(env, []Rule{{Op: OpDelete, Key: "A"}})
	if _, ok := out["A"]; ok {
		t.Error("expected A to be deleted")
	}
	if r := findResult(results, "A"); r == nil || !r.Applied {
		t.Error("expected applied result for A")
	}
}

func TestApply_Delete_Missing(t *testing.T) {
	env := map[string]string{"B": "2"}
	_, results := Apply(env, []Rule{{Op: OpDelete, Key: "MISSING"}})
	if r := findResult(results, "MISSING"); r == nil || r.Applied {
		t.Error("expected skipped result for MISSING")
	}
}

func TestApply_Rename(t *testing.T) {
	env := map[string]string{"OLD": "val"}
	out, results := Apply(env, []Rule{{Op: OpRename, Key: "OLD", To: "NEW"}})
	if out["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %s", out["NEW"])
	}
	if _, ok := out["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
	if r := findResult(results, "OLD"); r == nil || !r.Applied {
		t.Error("expected applied result for OLD")
	}
}

func TestSummary_Counts(t *testing.T) {
	results := []Result{
		{Applied: true},
		{Applied: true},
		{Applied: false},
	}
	applied, skipped := Summary(results)
	if applied != 2 || skipped != 1 {
		t.Errorf("expected 2/1, got %d/%d", applied, skipped)
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"set:FOO=bar", "delete:BAZ", "rename:OLD=NEW"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(rules))
	}
}

func TestParseRules_Invalid(t *testing.T) {
	_, err := ParseRules([]string{"badformat"})
	if err == nil {
		t.Error("expected error for bad format")
	}
}
