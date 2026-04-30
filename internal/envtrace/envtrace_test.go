package envtrace

import (
	"strings"
	"testing"
)

func findResult(entries []TraceEntry, key string) (TraceEntry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return TraceEntry{}, false
}

func TestApply_FileSource(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceFile, Origin: ".env", Data: map[string]string{"DB_HOST": "localhost"}},
		},
	}
	entries := Apply(opts)
	e, ok := findResult(entries, "DB_HOST")
	if !ok {
		t.Fatal("expected DB_HOST in results")
	}
	if e.Source != SourceFile {
		t.Errorf("expected source=file, got %s", e.Source)
	}
	if e.Origin != ".env" {
		t.Errorf("expected origin=.env, got %s", e.Origin)
	}
	if e.Value != "localhost" {
		t.Errorf("expected value=localhost, got %s", e.Value)
	}
}

func TestApply_FirstSourceWins(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceEnv, Origin: "process", Data: map[string]string{"PORT": "8080"}},
			{Kind: SourceFile, Origin: ".env", Data: map[string]string{"PORT": "3000"}},
		},
	}
	entries := Apply(opts)
	e, ok := findResult(entries, "PORT")
	if !ok {
		t.Fatal("expected PORT in results")
	}
	if e.Value != "8080" {
		t.Errorf("expected value=8080 (first source wins), got %s", e.Value)
	}
	if !e.Override {
		t.Error("expected Override=true when key exists in multiple sources")
	}
}

func TestApply_DefaultFallback(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceDefault, Origin: "defaults", Data: map[string]string{"LOG_LEVEL": "info"}},
		},
	}
	entries := Apply(opts)
	e, ok := findResult(entries, "LOG_LEVEL")
	if !ok {
		t.Fatal("expected LOG_LEVEL")
	}
	if e.Source != SourceDefault {
		t.Errorf("expected source=default, got %s", e.Source)
	}
}

func TestApply_NoOverride(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceFile, Origin: ".env", Data: map[string]string{"ONLY_HERE": "yes"}},
		},
	}
	entries := Apply(opts)
	e, _ := findResult(entries, "ONLY_HERE")
	if e.Override {
		t.Error("expected Override=false for key in single source")
	}
}

func TestFilterBySource(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceFile, Origin: ".env", Data: map[string]string{"A": "1"}},
			{Kind: SourceEnv, Origin: "process", Data: map[string]string{"B": "2"}},
		},
	}
	entries := Apply(opts)
	fileOnly := FilterBySource(entries, SourceFile)
	if len(fileOnly) != 1 || fileOnly[0].Key != "A" {
		t.Errorf("expected 1 file-source entry for A, got %+v", fileOnly)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	opts := Options{
		Sources: []SourceSpec{
			{Kind: SourceFile, Origin: ".env", Data: map[string]string{"X": "1", "Y": "2"}},
			{Kind: SourceEnv, Origin: "process", Data: map[string]string{"X": "override"}},
		},
	}
	entries := Apply(opts)
	summary := GetSummary(entries)
	if !strings.Contains(summary, "total=2") {
		t.Errorf("expected total=2 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "overridden=1") {
		t.Errorf("expected overridden=1 in summary, got: %s", summary)
	}
}
