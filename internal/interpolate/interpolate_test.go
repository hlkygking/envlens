package interpolate

import (
	"testing"
)

func findResult(results []Result, key string) *Result {
	for i := range results {
		if results[i].Key == key {
			return &results[i]
		}
	}
	return nil
}

func TestInterpolate_Basic(t *testing.T) {
	env := map[string]string{
		"BASE": "/app",
		"DATA": "${BASE}/data",
	}
	results := Interpolate(env)
	r := findResult(results, "DATA")
	if r == nil {
		t.Fatal("expected result for DATA")
	}
	if r.Resolved != "/app/data" {
		t.Errorf("expected /app/data, got %s", r.Resolved)
	}
}

func TestInterpolate_MissingRef(t *testing.T) {
	env := map[string]string{
		"PATH": "${MISSING}/bin",
	}
	results := Interpolate(env)
	r := findResult(results, "PATH")
	if r == nil {
		t.Fatal("expected result for PATH")
	}
	if len(r.Missing) == 0 {
		t.Error("expected missing refs")
	}
	if r.Missing[0] != "MISSING" {
		t.Errorf("expected MISSING, got %s", r.Missing[0])
	}
}

func TestInterpolate_NoRefs(t *testing.T) {
	env := map[string]string{"PLAIN": "hello"}
	results := Interpolate(env)
	r := findResult(results, "PLAIN")
	if r.Resolved != "hello" {
		t.Errorf("expected hello, got %s", r.Resolved)
	}
	if len(r.Refs) != 0 {
		t.Error("expected no refs")
	}
}

func TestInterpolate_SelfRef(t *testing.T) {
	env := map[string]string{"A": "${A}/extra"}
	results := Interpolate(env)
	r := findResult(results, "A")
	if r.Resolved != "${A}/extra" {
		t.Errorf("self-ref should not expand, got %s", r.Resolved)
	}
}

func TestGetSummary(t *testing.T) {
	env := map[string]string{
		"BASE": "/app",
		"DATA": "${BASE}/data",
		"BAD":  "${NOPE}/x",
	}
	results := Interpolate(env)
	s := GetSummary(results)
	if s.Total != 3 {
		t.Errorf("expected 3, got %d", s.Total)
	}
	if s.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", s.Missing)
	}
	if s.Resolved != 1 {
		t.Errorf("expected 1 resolved, got %d", s.Resolved)
	}
}

func TestValidate_Cycle(t *testing.T) {
	env := map[string]string{
		"A": "${B}/a",
		"B": "${A}/b",
	}
	cycles := Validate(env)
	if len(cycles) == 0 {
		t.Error("expected cycle detection")
	}
}
