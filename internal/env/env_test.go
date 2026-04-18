package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFromFile_Basic(t *testing.T) {
	p := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	s, err := FromFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if s.Vars["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", s.Vars["FOO"])
	}
	if s.Vars["BAZ"] != "qux" {
		t.Errorf("expected qux, got %s", s.Vars["BAZ"])
	}
}

func TestFromFile_IgnoresComments(t *testing.T) {
	p := writeTemp(t, "# comment\nKEY=val\n")
	s, err := FromFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Vars["# comment"]; ok {
		t.Error("comment should not be parsed as key")
	}
	if s.Vars["KEY"] != "val" {
		t.Errorf("expected val, got %s", s.Vars["KEY"])
	}
}

func TestFromFile_MissingFile(t *testing.T) {
	_, err := FromFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	s := FromMap("test", m)
	if s.Vars["A"] != "1" || s.Vars["B"] != "2" {
		t.Error("FromMap values mismatch")
	}
	// ensure copy
	m["A"] = "changed"
	if s.Vars["A"] != "1" {
		t.Error("FromMap should copy the map")
	}
}

func TestMerge(t *testing.T) {
	a := FromMap("a", map[string]string{"X": "1", "Y": "2"})
	b := FromMap("b", map[string]string{"Y": "override", "Z": "3"})
	merged := Merge(a, b)
	if merged.Vars["X"] != "1" {
		t.Errorf("expected X=1, got %s", merged.Vars["X"])
	}
	if merged.Vars["Y"] != "override" {
		t.Errorf("expected Y=override, got %s", merged.Vars["Y"])
	}
	if merged.Vars["Z"] != "3" {
		t.Errorf("expected Z=3, got %s", merged.Vars["Z"])
	}
}
