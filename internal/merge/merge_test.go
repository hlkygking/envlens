package merge

import (
	"testing"
)

func TestApply_NoConflict(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	over := map[string]string{"C": "3"}
	r, err := Apply(base, over, StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	if r.Merged["A"] != "1" || r.Merged["B"] != "2" || r.Merged["C"] != "3" {
		t.Errorf("unexpected merged map: %v", r.Merged)
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(r.Conflicts))
	}
	if len(r.Added) != 1 || r.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", r.Added)
	}
}

func TestApply_OverrideWins(t *testing.T) {
	base := map[string]string{"X": "old"}
	over := map[string]string{"X": "new"}
	r, err := Apply(base, over, StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	if r.Merged["X"] != "new" {
		t.Errorf("expected override value, got %q", r.Merged["X"])
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(r.Conflicts))
	}
}

func TestApply_BaseWins(t *testing.T) {
	base := map[string]string{"X": "old"}
	over := map[string]string{"X": "new"}
	r, err := Apply(base, over, StrategyBase)
	if err != nil {
		t.Fatal(err)
	}
	if r.Merged["X"] != "old" {
		t.Errorf("expected base value, got %q", r.Merged["X"])
	}
}

func TestApply_ErrorStrategy(t *testing.T) {
	base := map[string]string{"X": "old"}
	over := map[string]string{"X": "new"}
	_, err := Apply(base, over, StrategyError)
	if err == nil {
		t.Error("expected error on conflict with StrategyError")
	}
}

func TestApply_KeptKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	over := map[string]string{"B": "2"}
	r, err := Apply(base, over, StrategyOverride)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Kept) != 1 || r.Kept[0] != "A" {
		t.Errorf("expected Kept=[A], got %v", r.Kept)
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{
		Merged:    map[string]string{"A": "1", "B": "2"},
		Conflicts: []Conflict{{Key: "B"}},
		Added:     []string{"C"},
		Kept:      []string{"A"},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
