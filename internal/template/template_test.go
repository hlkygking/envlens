package template

import (
	"testing"
)

func findResult(results []Result, key string) (Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

func TestApply_Basic(t *testing.T) {
	tmpls := map[string]string{"GREETING": "Hello, {{NAME}}!"}
	env := map[string]string{"NAME": "World"}
	results := Apply(tmpls, env)
	r, ok := findResult(results, "GREETING")
	if !ok {
		t.Fatal("expected GREETING result")
	}
	if r.Rendered != "Hello, World!" {
		t.Errorf("unexpected rendered: %s", r.Rendered)
	}
	if !r.OK {
		t.Error("expected OK")
	}
}

func TestApply_MissingRef(t *testing.T) {
	tmpls := map[string]string{"URL": "http://{{HOST}}:{{PORT}}"}
	env := map[string]string{"HOST": "localhost"}
	results := Apply(tmpls, env)
	r, ok := findResult(results, "URL")
	if !ok {
		t.Fatal("expected URL result")
	}
	if r.OK {
		t.Error("expected not OK")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "PORT" {
		t.Errorf("expected PORT missing, got %v", r.Missing)
	}
}

func TestApply_NoPlaceholders(t *testing.T) {
	tmpls := map[string]string{"STATIC": "just-a-value"}
	results := Apply(tmpls, map[string]string{})
	r, ok := findResult(results, "STATIC")
	if !ok {
		t.Fatal("expected STATIC result")
	}
	if r.Rendered != "just-a-value" || !r.OK {
		t.Errorf("unexpected result: %+v", r)
	}
}

func TestToMap_OnlyOK(t *testing.T) {
	results := []Result{
		{Key: "A", Rendered: "val-a", OK: true},
		{Key: "B", Rendered: "<missing:X>", OK: false},
	}
	m := ToMap(results)
	if _, ok := m["A"]; !ok {
		t.Error("expected A in map")
	}
	if _, ok := m["B"]; ok {
		t.Error("expected B excluded from map")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	results := []Result{
		{OK: true}, {OK: true}, {OK: false},
	}
	ok, failed := GetSummary(results)
	if ok != 2 || failed != 1 {
		t.Errorf("expected 2 ok 1 failed, got %d %d", ok, failed)
	}
}
