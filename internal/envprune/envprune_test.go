package envprune_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/envprune"
)

func findResult(results []envprune.Result, key string) (envprune.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return envprune.Result{}, false
}

func TestApply_RemoveEmpty(t *testing.T) {
	env := map[string]string{"FOO": "bar", "EMPTY": "", "BAZ": "qux"}
	results := envprune.Apply(env, envprune.Options{RemoveEmpty: true})
	r, ok := findResult(results, "EMPTY")
	if !ok || !r.Pruned {
		t.Fatal("expected EMPTY to be pruned")
	}
	r2, _ := findResult(results, "FOO")
	if r2.Pruned {
		t.Fatal("FOO should not be pruned")
	}
}

func TestApply_RemovePrefixes(t *testing.T) {
	env := map[string]string{"DEBUG_VERBOSE": "1", "APP_NAME": "myapp", "DEBUG_LEVEL": "3"}
	results := envprune.Apply(env, envprune.Options{RemovePrefixes: []string{"DEBUG_"}})
	for _, key := range []string{"DEBUG_VERBOSE", "DEBUG_LEVEL"} {
		r, ok := findResult(results, key)
		if !ok || !r.Pruned {
			t.Errorf("expected %s to be pruned", key)
		}
	}
	r, _ := findResult(results, "APP_NAME")
	if r.Pruned {
		t.Fatal("APP_NAME should not be pruned")
	}
}

func TestApply_RemovePatterns(t *testing.T) {
	env := map[string]string{"TMP_FOO": "1", "TEMP_BAR": "2", "KEEP_ME": "3"}
	results := envprune.Apply(env, envprune.Options{RemovePatterns: []string{"^TMP_|^TEMP_"}})
	for _, key := range []string{"TMP_FOO", "TEMP_BAR"} {
		r, ok := findResult(results, key)
		if !ok || !r.Pruned {
			t.Errorf("expected %s to be pruned", key)
		}
	}
}

func TestApply_KeepOnly(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	results := envprune.Apply(env, envprune.Options{KeepOnly: []string{"APP_"}})
	r, _ := findResult(results, "DB_URL")
	if !r.Pruned {
		t.Fatal("DB_URL should be pruned (not in keep-only)")
	}
	r2, _ := findResult(results, "APP_HOST")
	if r2.Pruned {
		t.Fatal("APP_HOST should be retained")
	}
}

func TestToMap_ExcludesPruned(t *testing.T) {
	env := map[string]string{"A": "", "B": "val"}
	results := envprune.Apply(env, envprune.Options{RemoveEmpty: true})
	m := envprune.ToMap(results)
	if _, ok := m["A"]; ok {
		t.Fatal("pruned key A should not appear in map")
	}
	if m["B"] != "val" {
		t.Fatal("B should be in map")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"X": "", "Y": "", "Z": "keep"}
	results := envprune.Apply(env, envprune.Options{RemoveEmpty: true})
	pruned, retained := envprune.GetSummary(results)
	if pruned != 2 {
		t.Errorf("expected 2 pruned, got %d", pruned)
	}
	if retained != 1 {
		t.Errorf("expected 1 retained, got %d", retained)
	}
}
