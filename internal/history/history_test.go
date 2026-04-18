package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "history.json")
}

func TestAppendAndList(t *testing.T) {
	p := tempPath(t)
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := Append(p, "v1", env); err != nil {
		t.Fatalf("Append error: %v", err)
	}
	entries, err := List(p)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Label != "v1" {
		t.Errorf("expected label v1, got %s", entries[0].Label)
	}
	if entries[0].Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080")
	}
}

func TestAppend_MultipleEntries(t *testing.T) {
	p := tempPath(t)
	for _, label := range []string{"v1", "v2", "v3"} {
		if err := Append(p, label, map[string]string{"X": label}); err != nil {
			t.Fatalf("Append error: %v", err)
		}
	}
	entries, _ := List(p)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestLatest_ReturnsLast(t *testing.T) {
	p := tempPath(t)
	_ = Append(p, "first", map[string]string{})
	_ = Append(p, "last", map[string]string{"K": "V"})
	e, err := Latest(p)
	if err != nil {
		t.Fatalf("Latest error: %v", err)
	}
	if e.Label != "last" {
		t.Errorf("expected label 'last', got %s", e.Label)
	}
}

func TestLatest_EmptyHistory(t *testing.T) {
	p := tempPath(t)
	_, err := Latest(p)
	if err == nil {
		t.Error("expected error for empty history")
	}
}

func TestList_MissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nonexistent.json")
	entries, err := List(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := List(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
