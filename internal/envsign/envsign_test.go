package envsign

import (
	"testing"
)

func findEntry(entries []Entry, key string) (Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return Entry{}, false
}

func TestSign_ProducesSignatures(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	entries := Sign(env, "mysecret")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Signature == "" {
			t.Errorf("expected non-empty signature for %s", e.Key)
		}
		if !e.Valid {
			t.Errorf("expected Valid=true for signed entry %s", e.Key)
		}
	}
}

func TestVerify_AllValid(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	signed := Sign(env, "secret123")
	sigs := make(map[string]string)
	for _, e := range signed {
		sigs[e.Key] = e.Signature
	}

	results := Verify(env, sigs, "secret123")
	for _, r := range results {
		if !r.Valid {
			t.Errorf("expected Valid=true for key %s", r.Key)
		}
	}
}

func TestVerify_TamperedValue(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "original"}
	signed := Sign(env, "mykey")
	sigs := map[string]string{signed[0].Key: signed[0].Signature}

	tampered := map[string]string{"SECRET_KEY": "tampered"}
	results := Verify(tampered, sigs, "mykey")
	e, ok := findEntry(results, "SECRET_KEY")
	if !ok {
		t.Fatal("expected entry for SECRET_KEY")
	}
	if e.Valid {
		t.Error("expected Valid=false for tampered value")
	}
}

func TestVerify_MissingSignature(t *testing.T) {
	env := map[string]string{"NEW_KEY": "value"}
	results := Verify(env, map[string]string{}, "secret")
	e, ok := findEntry(results, "NEW_KEY")
	if !ok {
		t.Fatal("expected entry for NEW_KEY")
	}
	if e.Valid {
		t.Error("expected Valid=false when signature is missing")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	entries := []Entry{
		{Key: "A", Valid: true},
		{Key: "B", Valid: true},
		{Key: "C", Valid: false},
	}
	s := GetSummary(entries)
	if s.Total != 3 || s.Valid != 2 || s.Invalid != 1 {
		t.Errorf("unexpected summary: %+v", s)
	}
}
