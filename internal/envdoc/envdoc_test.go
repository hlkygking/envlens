package envdoc

import (
	"os"
	"testing"
)

func findResult(results []Result, key string) (Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

var sampleDocs = []Entry{
	{Key: "APP_PORT", Description: "HTTP port", Example: "8080", Required: true, Default: "3000"},
	{Key: "DB_URL", Description: "Database connection URL", Required: true},
	{Key: "LOG_LEVEL", Description: "Logging verbosity", Default: "info"},
}

func TestApply_Documented(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080", "DB_URL": "postgres://localhost"}
	results := Apply(env, sampleDocs)
	r, ok := findResult(results, "APP_PORT")
	if !ok || !r.Found {
		t.Fatal("expected APP_PORT to be documented")
	}
	if r.Description != "HTTP port" {
		t.Errorf("expected description 'HTTP port', got %q", r.Description)
	}
	if r.Default != "3000" {
		t.Errorf("expected default '3000', got %q", r.Default)
	}
}

func TestApply_Undocumented(t *testing.T) {
	env := map[string]string{"UNKNOWN_KEY": "value"}
	results := Apply(env, sampleDocs)
	r, ok := findResult(results, "UNKNOWN_KEY")
	if !ok {
		t.Fatal("expected UNKNOWN_KEY in results")
	}
	if r.Found {
		t.Error("expected UNKNOWN_KEY to be undocumented")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Found: true},
		{Key: "B", Found: true},
		{Key: "C", Found: false},
	}
	s := GetSummary(results)
	if s.Total != 3 || s.Documented != 2 || s.Undocumented != 1 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestParseDocFile_Basic(t *testing.T) {
	f, _ := os.CreateTemp("", "envdoc*.txt")
	defer os.Remove(f.Name())
	f.WriteString("# comment\n")
	f.WriteString("APP_PORT:HTTP port:8080:true:3000\n")
	f.WriteString("LOG_LEVEL:Verbosity::false:info\n")
	f.Close()

	entries, err := ParseDocFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "APP_PORT" || !entries[0].Required || entries[0].Default != "3000" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseDocFile_MissingFile(t *testing.T) {
	_, err := ParseDocFile("/nonexistent/path.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
