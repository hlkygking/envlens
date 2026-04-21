package envpivot_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/envpivot"
)

func findEntry(entries []envpivot.Entry, key string) (envpivot.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return envpivot.Entry{}, false
}

func TestApply_Basic(t *testing.T) {
	scopes := map[string]map[string]string{
		"dev":  {"PORT": "3000", "DEBUG": "true"},
		"prod": {"PORT": "8080", "DEBUG": "false"},
	}
	entries := envpivot.Apply(scopes)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	e, ok := findEntry(entries, "PORT")
	if !ok {
		t.Fatal("expected PORT entry")
	}
	if e.Values["dev"] != "3000" || e.Values["prod"] != "8080" {
		t.Errorf("unexpected PORT values: %v", e.Values)
	}
	if e.Uniform {
		t.Error("PORT should not be uniform")
	}
}

func TestApply_UniformKey(t *testing.T) {
	scopes := map[string]map[string]string{
		"dev":     {"APP": "myapp"},
		"staging": {"APP": "myapp"},
		"prod":    {"APP": "myapp"},
	}
	entries := envpivot.Apply(scopes)
	e, ok := findEntry(entries, "APP")
	if !ok {
		t.Fatal("expected APP entry")
	}
	if !e.Uniform {
		t.Error("APP should be uniform")
	}
	if len(e.Missing) != 0 {
		t.Errorf("expected no missing scopes, got %v", e.Missing)
	}
}

func TestApply_MissingInScope(t *testing.T) {
	scopes := map[string]map[string]string{
		"dev":  {"SECRET": "abc"},
		"prod": {},
	}
	entries := envpivot.Apply(scopes)
	e, ok := findEntry(entries, "SECRET")
	if !ok {
		t.Fatal("expected SECRET entry")
	}
	if len(e.Missing) != 1 || e.Missing[0] != "prod" {
		t.Errorf("expected prod missing, got %v", e.Missing)
	}
	if e.Uniform {
		t.Error("SECRET should not be uniform when missing in a scope")
	}
}

func TestApply_EmptyScopes(t *testing.T) {
	entries := envpivot.Apply(nil)
	if entries != nil {
		t.Error("expected nil for empty input")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	entries := []envpivot.Entry{
		{Key: "A", Uniform: true, Missing: nil},
		{Key: "B", Uniform: false, Missing: []string{"prod"}},
		{Key: "C", Uniform: false, Missing: nil},
	}
	s := envpivot.GetSummary(entries)
	if s.TotalKeys != 3 {
		t.Errorf("expected 3 total, got %d", s.TotalKeys)
	}
	if s.UniformKeys != 1 {
		t.Errorf("expected 1 uniform, got %d", s.UniformKeys)
	}
	if s.DivergentKeys != 2 {
		t.Errorf("expected 2 divergent, got %d", s.DivergentKeys)
	}
	if s.MissingInAny != 1 {
		t.Errorf("expected 1 missing-in-any, got %d", s.MissingInAny)
	}
}
