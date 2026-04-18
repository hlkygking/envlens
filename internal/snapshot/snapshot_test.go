package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}
	path := tempPath(t)
	if err := Save(path, "test-label", env); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if s.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", s.Label)
	}
	if s.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", s.Env["APP_ENV"])
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	path := tempPath(t)
	if err := Save(path, "label", map[string]string{"K": "V"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after Save")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Error("expected error loading invalid JSON")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_TimestampIsRecent(t *testing.T) {
	before := time.Now().UTC()
	path := tempPath(t)
	_ = Save(path, "ts-test", map[string]string{})
	s, _ := Load(path)
	after := time.Now().UTC()
	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v not in expected range [%v, %v]", s.CreatedAt, before, after)
	}
}
