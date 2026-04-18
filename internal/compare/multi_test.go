package compare

import (
	"testing"
)

func TestMultiFiles_Basic(t *testing.T) {
	base := writeTemp(t, "A=1\nB=2\n")
	t1 := writeTemp(t, "A=1\nB=changed\n")
	t2 := writeTemp(t, "A=new\nC=3\n")
	res, err := MultiFiles(base, []string{t1, t2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Base != base {
		t.Errorf("base path mismatch")
	}
	if len(res.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(res.Targets))
	}
}

func TestMultiFiles_AllKeys(t *testing.T) {
	base := writeTemp(t, "A=1\n")
	t1 := writeTemp(t, "B=2\n")
	t2 := writeTemp(t, "C=3\n")
	res, err := MultiFiles(base, []string{t1, t2})
	if err != nil {
		t.Fatal(err)
	}
	keys := res.AllKeys()
	if len(keys) < 2 {
		t.Errorf("expected at least 2 unique keys, got %d", len(keys))
	}
}

func TestMultiFiles_MissingTarget(t *testing.T) {
	base := writeTemp(t, "A=1\n")
	_, err := MultiFiles(base, []string{"/no/such/file.env"})
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestMultiFiles_EmptyTargets(t *testing.T) {
	base := writeTemp(t, "A=1\n")
	res, err := MultiFiles(base, []string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Targets) != 0 {
		t.Error("expected no targets")
	}
}
