package compare

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envlens-*.env")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestFiles_Basic(t *testing.T) {
	base := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	target := writeTemp(t, "FOO=changed\nNEW=val\n")
	res, err := Files(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BaseFile != base || res.TargetFile != target {
		t.Error("file paths not stored correctly")
	}
	if len(res.Entries) == 0 {
		t.Error("expected entries")
	}
}

func TestFiles_Summary(t *testing.T) {
	base := writeTemp(t, "A=1\nB=2\n")
	target := writeTemp(t, "A=1\nC=3\n")
	res, err := Files(base, target)
	if err != nil {
		t.Fatal(err)
	}
	s := res.Summary()
	if s["removed"] != 1 {
		t.Errorf("expected 1 removed, got %d", s["removed"])
	}
	if s["added"] != 1 {
		t.Errorf("expected 1 added, got %d", s["added"])
	}
}

func TestFiles_MissingBase(t *testing.T) {
	_, err := Files("/nonexistent/base.env", "/nonexistent/target.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFiles_SortedKeys(t *testing.T) {
	base := writeTemp(t, "Z=1\nA=2\nM=3\n")
	target := writeTemp(t, "Z=1\nA=9\nM=3\n")
	res, err := Files(base, target)
	if err != nil {
		t.Fatal(err)
	}
	if res.Entries[0].Key != "A" {
		t.Errorf("expected sorted first key A, got %s", res.Entries[0].Key)
	}
}
