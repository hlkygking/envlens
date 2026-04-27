package envexpand_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/envexpand"
)

func findResult(results []envexpand.Result, key string) (envexpand.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return envexpand.Result{}, false
}

func TestApply_Unchanged(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	results := envexpand.Apply(env, false)
	r, ok := findResult(results, "FOO")
	if !ok {
		t.Fatal("expected FOO in results")
	}
	if r.Status != "unchanged" {
		t.Errorf("expected unchanged, got %s", r.Status)
	}
	if r.Expanded != "bar" {
		t.Errorf("expected bar, got %s", r.Expanded)
	}
}

func TestApply_ExpandsBraceRef(t *testing.T) {
	env := map[string]string{"BASE": "/app", "LOG": "${BASE}/logs"}
	results := envexpand.Apply(env, false)
	r, ok := findResult(results, "LOG")
	if !ok {
		t.Fatal("expected LOG in results")
	}
	if r.Status != "ok" {
		t.Errorf("expected ok, got %s", r.Status)
	}
	if r.Expanded != "/app/logs" {
		t.Errorf("expected /app/logs, got %s", r.Expanded)
	}
}

func TestApply_UnresolvedRef(t *testing.T) {
	env := map[string]string{"URL": "http://${HOST}:8080"}
	results := envexpand.Apply(env, false)
	r, ok := findResult(results, "URL")
	if !ok {
		t.Fatal("expected URL in results")
	}
	if r.Status != "unresolved" {
		t.Errorf("expected unresolved, got %s", r.Status)
	}
}

func TestApply_BareRef(t *testing.T) {
	env := map[string]string{"HOME_DIR": "/home/user", "CONF": "$HOME_DIR/.config"}
	results := envexpand.Apply(env, false)
	r, ok := findResult(results, "CONF")
	if !ok {
		t.Fatal("expected CONF")
	}
	if r.Expanded != "/home/user/.config" {
		t.Errorf("unexpected expansion: %s", r.Expanded)
	}
}

func TestToMap_OnlyResolved(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${MISSING}_val",
	}
	results := envexpand.Apply(env, false)
	m := envexpand.ToMap(results)
	if _, ok := m["A"]; !ok {
		t.Error("expected A in map")
	}
	if _, ok := m["B"]; ok {
		t.Error("expected B excluded from map due to unresolved status")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"PLAIN":     "value",
		"EXPANDED":  "${PLAIN}_suffix",
		"BROKEN":    "${NOPE}",
	}
	results := envexpand.Apply(env, false)
	s := envexpand.GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected unchanged 1, got %d", s.Unchanged)
	}
	if s.Expanded != 1 {
		t.Errorf("expected expanded 1, got %d", s.Expanded)
	}
	if s.Unresolved != 1 {
		t.Errorf("expected unresolved 1, got %d", s.Unresolved)
	}
}
