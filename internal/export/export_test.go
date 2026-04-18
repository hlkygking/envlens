package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"SECRET":   "abc123",
}

func TestRender_Dotenv(t *testing.T) {
	out, err := Render(sampleEnv, FormatDotenv)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "APP_ENV=") {
		t.Error("expected APP_ENV in dotenv output")
	}
	if !strings.Contains(out, "DB_HOST=") {
		t.Error("expected DB_HOST in dotenv output")
	}
}

func TestRender_Shell(t *testing.T) {
	out, err := Render(sampleEnv, FormatShell)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Error("expected export prefix in shell output")
	}
}

func TestRender_JSON(t *testing.T) {
	out, err := Render(sampleEnv, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", m["APP_ENV"])
	}
}

func TestRender_UnknownFormat(t *testing.T) {
	_, err := Render(sampleEnv, Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExport_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	if err := Export(sampleEnv, FormatDotenv, path); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "APP_ENV=") {
		t.Error("expected APP_ENV in written file")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	out, _ := Render(sampleEnv, FormatDotenv)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// APP_ENV < DB_HOST < SECRET alphabetically
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected first line APP_ENV, got %s", lines[0])
	}
}
