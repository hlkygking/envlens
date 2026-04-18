package promote

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

func TestApply_Promoted(t *testing.T) {
	src := map[string]string{"NEW_KEY": "hello"}
	dst := map[string]string{}
	results, err := Apply(src, dst, StrategyKeep)
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(results, "NEW_KEY")
	if r == nil || r.Status != "promoted" {
		t.Errorf("expected promoted, got %+v", r)
	}
	if dst["NEW_KEY"] != "hello" {
		t.Errorf("expected dst to contain NEW_KEY")
	}
}

func TestApply_OverrideWins(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}
	results, err := Apply(src, dst, StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(results, "KEY")
	if r == nil || r.Status != "promoted" || !r.Conflict {
		t.Errorf("expected promoted conflict, got %+v", r)
	}
	if dst["KEY"] != "new" {
		t.Errorf("expected dst KEY=new")
	}
}

func TestApply_KeepWins(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}
	results, err := Apply(src, dst, StrategyKeep)
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(results, "KEY")
	if r == nil || r.Status != "skipped" {
		t.Errorf("expected skipped, got %+v", r)
	}
	if dst["KEY"] != "old" {
		t.Errorf("expected dst KEY=old")
	}
}

func TestApply_ErrorStrategy(t *testing.T) {
	src := map[string]string{"KEY": "new"}
	dst := map[string]string{"KEY": "old"}
	_, err := Apply(src, dst, StrategyError)
	if err == nil {
		t.Error("expected error for conflict")
	}
}

func TestApply_IdenticalSkipped(t *testing.T) {
	src := map[string]string{"KEY": "same"}
	dst := map[string]string{"KEY": "same"}
	results, err := Apply(src, dst, StrategyError)
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(results, "KEY")
	if r == nil || r.Status != "skipped" {
		t.Errorf("expected skipped for identical, got %+v", r)
	}
}

func TestGetSummary(t *testing.T) {
	results := []Result{
		{Key: "A", Status: "promoted"},
		{Key: "B", Status: "promoted", Conflict: true},
		{Key: "C", Status: "skipped"},
	}
	s := GetSummary(results)
	if s.Promoted != 2 {
		t.Errorf("expected 2 promoted, got %d", s.Promoted)
	}
	if s.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", s.Conflicts)
	}
	if s.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", s.Skipped)
	}
}
