package envnorm

import (
	"testing"
)

func findResult(results []Result, orig string) *Result {
	for i := range results {
		if results[i].OriginalKey == orig {
			return &results[i]
		}
	}
	return nil
}

func TestApply_Uppercase(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	res := Apply(env, []NormMode{ModeUppercase})
	r := findResult(res, "db_host")
	if r == nil {
		t.Fatal("result not found")
	}
	if r.NormalizedKey != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", r.NormalizedKey)
	}
	if !r.Changed {
		t.Error("expected Changed=true")
	}
}

func TestApply_Lowercase(t *testing.T) {
	env := map[string]string{"APP_NAME": "envlens"}
	res := Apply(env, []NormMode{ModeLowercase})
	r := findResult(res, "APP_NAME")
	if r.NormalizedKey != "app_name" {
		t.Errorf("expected app_name, got %s", r.NormalizedKey)
	}
}

func TestApply_SnakeCase(t *testing.T) {
	env := map[string]string{"DbHost": "localhost"}
	res := Apply(env, []NormMode{ModeSnakeCase})
	r := findResult(res, "DbHost")
	if r.NormalizedKey != "db_host" {
		t.Errorf("expected db_host, got %s", r.NormalizedKey)
	}
}

func TestApply_StripSpace(t *testing.T) {
	env := map[string]string{"DB HOST": "localhost"}
	res := Apply(env, []NormMode{ModeStripSpace})
	r := findResult(res, "DB HOST")
	if r.NormalizedKey != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", r.NormalizedKey)
	}
}

func TestApply_Unchanged(t *testing.T) {
	env := map[string]string{"ALREADY_UPPER": "val"}
	res := Apply(env, []NormMode{ModeUppercase})
	r := findResult(res, "ALREADY_UPPER")
	if r.Changed {
		t.Error("expected Changed=false")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "app_port": "8080"}
	res := Apply(env, []NormMode{ModeUppercase})
	m := ToMap(res)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", m["DB_HOST"])
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"db_host": "a", "ALREADY": "b"}
	res := Apply(env, []NormMode{ModeUppercase})
	s := GetSummary(res)
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
	if s.Changed != 1 {
		t.Errorf("expected changed 1, got %d", s.Changed)
	}
}
