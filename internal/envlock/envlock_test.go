package envlock

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func findResult(results []CheckResult, key string) (CheckResult, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return CheckResult{}, false
}

func TestLock_ProducesEntries(t *testing.T) {
	env := map[string]string{"DB_URL": "postgres://localhost", "PORT": "5432"}
	lf := Lock(env)
	if len(lf.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lf.Entries))
	}
	if lf.Version != 1 {
		t.Errorf("expected version 1, got %d", lf.Version)
	}
	for _, e := range lf.Entries {
		if e.Hash == "" {
			t.Errorf("expected non-empty hash for key %s", e.Key)
		}
		if e.LockedAt == "" {
			t.Errorf("expected non-empty locked_at for key %s", e.Key)
		}
	}
}

func TestLock_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	lf := Lock(env)
	if lf.Entries[0].Key != "A_KEY" || lf.Entries[1].Key != "M_KEY" || lf.Entries[2].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", lf.Entries)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"SECRET": "abc123", "HOST": "localhost"}
	lf := Lock(env)
	path := filepath.Join(t.TempDir(), "env.lock")
	if err := Save(path, lf); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Entries) != len(lf.Entries) {
		t.Errorf("entry count mismatch: got %d, want %d", len(loaded.Entries), len(lf.Entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.lock")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCheck_AllOK(t *testing.T) {
	env := map[string]string{"API_KEY": "secret", "PORT": "8080"}
	lf := Lock(env)
	results := Check(lf, env)
	for _, r := range results {
		if r.Status != "ok" {
			t.Errorf("expected ok for %s, got %s", r.Key, r.Status)
		}
	}
}

func TestCheck_Drifted(t *testing.T) {
	env := map[string]string{"API_KEY": "original"}
	lf := Lock(env)
	modified := map[string]string{"API_KEY": "changed"}
	results := Check(lf, modified)
	r, ok := findResult(results, "API_KEY")
	if !ok {
		t.Fatal("API_KEY not found in results")
	}
	if r.Status != "drifted" {
		t.Errorf("expected drifted, got %s", r.Status)
	}
}

func TestCheck_Missing(t *testing.T) {
	env := map[string]string{"DB_PASS": "hunter2"}
	lf := Lock(env)
	results := Check(lf, map[string]string{})
	r, ok := findResult(results, "DB_PASS")
	if !ok {
		t.Fatal("DB_PASS not found")
	}
	if r.Status != "missing" {
		t.Errorf("expected missing, got %s", r.Status)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []CheckResult{
		{Key: "A", Status: "ok"},
		{Key: "B", Status: "drifted"},
		{Key: "C", Status: "missing"},
		{Key: "D", Status: "ok"},
	}
	ok, drifted, missing := GetSummary(results)
	if ok != 2 || drifted != 1 || missing != 1 {
		t.Errorf("unexpected counts: ok=%d drifted=%d missing=%d", ok, drifted, missing)
	}
}

func TestSave_WritesValidJSON(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	lf := Lock(env)
	path := filepath.Join(t.TempDir(), "env.lock")
	Save(path, lf)
	data, _ := os.ReadFile(path)
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Errorf("saved file is not valid JSON: %v", err)
	}
}
