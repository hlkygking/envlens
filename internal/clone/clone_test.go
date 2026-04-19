package clone_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your/envlens/internal/clone"
	"github.com/your/envlens/internal/env"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestApply_Basic(t *testing.T) {
	src := writeTemp(t, "APP=hello\nDB_HOST=localhost\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := clone.Apply(src, dst, clone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", res.Keys)
	}
	m, _ := env.FromFile(dst)
	if m["APP"] != "hello" {
		t.Errorf("expected APP=hello, got %s", m["APP"])
	}
}

func TestApply_RedactSensitive(t *testing.T) {
	src := writeTemp(t, "APP=hello\nAPI_KEY=supersecret\nDB_PASSWORD=abc123\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	res, err := clone.Apply(src, dst, clone.Options{RedactSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Redacted {
		t.Error("expected Redacted=true")
	}
	m, _ := env.FromFile(dst)
	if m["API_KEY"] != "" {
		t.Errorf("expected API_KEY to be redacted, got %s", m["API_KEY"])
	}
	if m["APP"] != "hello" {
		t.Errorf("expected APP=hello, got %s", m["APP"])
	}
}

func TestApply_ExtraKeys(t *testing.T) {
	src := writeTemp(t, "APP=hello\n")
	dst := filepath.Join(t.TempDir(), "out.env")

	_, err := clone.Apply(src, dst, clone.Options{Extra: map[string]string{"EXTRA_KEY": "injected"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, _ := env.FromFile(dst)
	if m["EXTRA_KEY"] != "injected" {
		t.Errorf("expected EXTRA_KEY=injected, got %s", m["EXTRA_KEY"])
	}
}

func TestApply_MissingSource(t *testing.T) {
	_, err := clone.Apply("/nonexistent.env", "/tmp/out.env", clone.Options{})
	if err == nil {
		t.Error("expected error for missing source")
	}
	if !strings.Contains(err.Error(), "clone") {
		t.Errorf("expected clone prefix in error, got: %v", err)
	}
}
