package resolve

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "envlens",
		"LOG_LEVEL": "",
		"PORT":     "8080",
	}
}

func findResult(results []Result, key string) (Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

func TestResolve_FileSource(t *testing.T) {
	results := Resolve(baseEnv(), Options{})
	r, ok := findResult(results, "APP_NAME")
	if !ok {
		t.Fatal("APP_NAME not found")
	}
	if r.Source != "file" || r.Value != "envlens" {
		t.Errorf("unexpected: %+v", r)
	}
}

func TestResolve_DefaultFallback(t *testing.T) {
	opts := Options{
		Defaults: map[string]string{"LOG_LEVEL": "info"},
	}
	results := Resolve(baseEnv(), opts)
	r, ok := findResult(results, "LOG_LEVEL")
	if !ok {
		t.Fatal("LOG_LEVEL not found")
	}
	if r.Source != "default" || r.Value != "info" {
		t.Errorf("unexpected: %+v", r)
	}
}

func TestResolve_MissingKey(t *testing.T) {
	results := Resolve(baseEnv(), Options{})
	r, ok := findResult(results, "LOG_LEVEL")
	if !ok {
		t.Fatal("LOG_LEVEL not found")
	}
	if r.Source != "missing" || r.Resolved {
		t.Errorf("expected missing: %+v", r)
	}
}

func TestResolve_DefaultOnlyKey(t *testing.T) {
	opts := Options{
		Defaults: map[string]string{"TIMEOUT": "30s"},
	}
	results := Resolve(map[string]string{}, opts)
	r, ok := findResult(results, "TIMEOUT")
	if !ok {
		t.Fatal("TIMEOUT not found")
	}
	if r.Source != "default" || r.Value != "30s" {
		t.Errorf("unexpected: %+v", r)
	}
}

func TestMissingKeys(t *testing.T) {
	results := Resolve(baseEnv(), Options{})
	missing := MissingKeys(results)
	if len(missing) != 1 || missing[0] != "LOG_LEVEL" {
		t.Errorf("expected [LOG_LEVEL], got %v", missing)
	}
}

func TestSummary(t *testing.T) {
	results := Resolve(baseEnv(), Options{})
	s := Summary(results)
	if s == "" {
		t.Error("empty summary")
	}
}

func TestToMap_OnlyResolved(t *testing.T) {
	results := Resolve(baseEnv(), Options{})
	m := ToMap(results)
	if _, ok := m["LOG_LEVEL"]; ok {
		t.Error("missing key should not appear in map")
	}
	if m["APP_NAME"] != "envlens" {
		t.Errorf("unexpected APP_NAME: %s", m["APP_NAME"])
	}
}

func TestResolve_EmptyEnvAndNoDefaults(t *testing.T) {
	results := Resolve(map[string]string{}, Options{})
	if len(results) != 0 {
		t.Errorf("expected no results for empty env, got %d", len(results))
	}
}
