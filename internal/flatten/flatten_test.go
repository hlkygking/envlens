package flatten

import (
	"testing"
)

func findResult(results []Result, original string) (Result, bool) {
	for _, r := range results {
		if r.OriginalKey == original {
			return r, true
		}
	}
	return Result{}, false
}

func TestApply_ReplaceDots(t *testing.T) {
	env := map[string]string{"app.name": "envlens"}
	results := Apply(env, Options{})
	r, ok := findResult(results, "app.name")
	if !ok {
		t.Fatal("expected result for app.name")
	}
	if r.FlatKey != "app_name" {
		t.Errorf("expected app_name, got %s", r.FlatKey)
	}
	if !r.Changed {
		t.Error("expected Changed=true")
	}
}

func TestApply_ReplaceDashes(t *testing.T) {
	env := map[string]string{"my-key": "val"}
	results := Apply(env, Options{Separator: "_"})
	r, ok := findResult(results, "my-key")
	if !ok {
		t.Fatal("expected result for my-key")
	}
	if r.FlatKey != "my_key" {
		t.Errorf("expected my_key, got %s", r.FlatKey)
	}
}

func TestApply_Uppercase(t *testing.T) {
	env := map[string]string{"db.host": "localhost"}
	results := Apply(env, Options{Uppercase: true})
	r, ok := findResult(results, "db.host")
	if !ok {
		t.Fatal("expected result")
	}
	if r.FlatKey != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", r.FlatKey)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	env := map[string]string{"APP_NAME": "x", "APP_PORT": "8080"}
	results := Apply(env, Options{StripPrefix: "APP_"})
	r, ok := findResult(results, "APP_NAME")
	if !ok {
		t.Fatal("expected result")
	}
	if r.FlatKey != "NAME" {
		t.Errorf("expected NAME, got %s", r.FlatKey)
	}
}

func TestApply_NoChange(t *testing.T) {
	env := map[string]string{"PLAIN_KEY": "value"}
	results := Apply(env, Options{})
	r, ok := findResult(results, "PLAIN_KEY")
	if !ok {
		t.Fatal("expected result")
	}
	if r.Changed {
		t.Error("expected Changed=false for already flat key")
	}
}

func TestToMap_RoundTrip(t *testing.T) {
	env := map[string]string{"a.b": "1", "c.d": "2"}
	m := ToMap(Apply(env, Options{}))
	if m["a_b"] != "1" || m["c_d"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{"a.b": "1", "PLAIN": "2"}
	s := GetSummary(Apply(env, Options{}))
	if s.Total != 2 {
		t.Errorf("expected Total=2, got %d", s.Total)
	}
	if s.Changed != 1 {
		t.Errorf("expected Changed=1, got %d", s.Changed)
	}
}
