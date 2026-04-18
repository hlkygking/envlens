package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "envlens")
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = "."
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, b)
	}
	return out
}

func TestMain_TextDiff(t *testing.T) {
	bin := buildBinary(t)
	base := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	target := writeTemp(t, "FOO=changed\nNEW=value\n")

	out, err := exec.Command(bin, base, target).Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := string(out)
	if !strings.Contains(body, "MODIFIED") {
		t.Errorf("expected MODIFIED section, got:\n%s", body)
	}
	if !strings.Contains(body, "ADDED") {
		t.Errorf("expected ADDED section, got:\n%s", body)
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	base := writeTemp(t, "FOO=bar\n")
	target := writeTemp(t, "FOO=bar\nEXTRA=1\n")

	out, err := exec.Command(bin, "-format=json", base, target).Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "\"key\"") {
		t.Errorf("expected JSON keys, got:\n%s", string(out))
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit for missing args")
	}
}
