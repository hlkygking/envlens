package envset

import (
	"testing"
)

func baseMap() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestNew_StoresEntries(t *testing.T) {
	s := New("test", baseMap())
	if s.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", s.Len())
	}
	if s.Name != "test" {
		t.Errorf("expected name 'test', got %s", s.Name)
	}
}

func TestGet_ExistingKey(t *testing.T) {
	s := New("s", baseMap())
	v, ok := s.Get("APP_HOST")
	if !ok || v != "localhost" {
		t.Errorf("expected localhost, got %q ok=%v", v, ok)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := New("s", baseMap())
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestSet_AddsKey(t *testing.T) {
	s := New("s", baseMap())
	s.Set("NEW_KEY", "value")
	if !s.Contains("NEW_KEY") {
		t.Error("expected NEW_KEY to be present")
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := New("s", baseMap())
	s.Delete("APP_PORT")
	if s.Contains("APP_PORT") {
		t.Error("expected APP_PORT to be deleted")
	}
	if s.Len() != 2 {
		t.Errorf("expected 2 entries after delete, got %d", s.Len())
	}
}

func TestKeys_Sorted(t *testing.T) {
	s := New("s", baseMap())
	keys := s.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestToMap_IsACopy(t *testing.T) {
	s := New("s", baseMap())
	m := s.ToMap()
	m["INJECTED"] = "x"
	if s.Contains("INJECTED") {
		t.Error("ToMap should return a copy, not a reference")
	}
}

func TestIntersect(t *testing.T) {
	a := New("a", map[string]string{"X": "1", "Y": "2"})
	b := New("b", map[string]string{"Y": "2", "Z": "3"})
	common := Intersect(a, b)
	if len(common) != 1 || common[0] != "Y" {
		t.Errorf("expected [Y], got %v", common)
	}
}

func TestUnion(t *testing.T) {
	a := New("a", map[string]string{"X": "1"})
	b := New("b", map[string]string{"Y": "2"})
	all := Union(a, b)
	if len(all) != 2 {
		t.Errorf("expected 2 keys, got %v", all)
	}
}

func TestSubtract(t *testing.T) {
	a := New("a", map[string]string{"X": "1", "Y": "2"})
	b := New("b", map[string]string{"X": "1"})
	only := Subtract(a, b)
	if len(only) != 1 || only[0] != "Y" {
		t.Errorf("expected [Y], got %v", only)
	}
}

func TestHasPrefix(t *testing.T) {
	s := New("s", baseMap())
	matches := s.HasPrefix("APP_")
	if len(matches) != 2 {
		t.Errorf("expected 2 APP_ keys, got %v", matches)
	}
}
