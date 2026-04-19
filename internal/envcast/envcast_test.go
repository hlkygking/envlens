package envcast

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

func TestApply_String(t *testing.T) {
	env := map[string]string{"NAME": "alice"}
	results := Apply(env, []Rule{{Key: "NAME", Type: TypeString}})
	r, ok := findResult(results, "NAME")
	if !ok || !r.OK || r.CastValue != "alice" {
		t.Errorf("expected string cast to succeed, got %+v", r)
	}
}

func TestApply_Int(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	results := Apply(env, []Rule{{Key: "PORT", Type: TypeInt}})
	r, ok := findResult(results, "PORT")
	if !ok || !r.OK || r.CastValue != 8080 {
		t.Errorf("expected int cast to 8080, got %+v", r)
	}
}

func TestApply_Float(t *testing.T) {
	env := map[string]string{"RATIO": "3.14"}
	results := Apply(env, []Rule{{Key: "RATIO", Type: TypeFloat}})
	r, ok := findResult(results, "RATIO")
	if !ok || !r.OK {
		t.Errorf("expected float cast to succeed, got %+v", r)
	}
}

func TestApply_Bool(t *testing.T) {
	env := map[string]string{"DEBUG": "true"}
	results := Apply(env, []Rule{{Key: "DEBUG", Type: TypeBool}})
	r, ok := findResult(results, "DEBUG")
	if !ok || !r.OK || r.CastValue != true {
		t.Errorf("expected bool cast to true, got %+v", r)
	}
}

func TestApply_InvalidInt(t *testing.T) {
	env := map[string]string{"PORT": "abc"}
	results := Apply(env, []Rule{{Key: "PORT", Type: TypeInt}})
	r, _ := findResult(results, "PORT")
	if r.OK || r.Error == "" {
		t.Errorf("expected int cast to fail, got %+v", r)
	}
}

func TestApply_MissingKey(t *testing.T) {
	env := map[string]string{}
	results := Apply(env, []Rule{{Key: "MISSING", Type: TypeString}})
	r, _ := findResult(results, "MISSING")
	if r.OK || r.Error != "key not found" {
		t.Errorf("expected missing key error, got %+v", r)
	}
}

func TestSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "1", "B": "bad"}
	rules := []Rule{{Key: "A", Type: TypeInt}, {Key: "B", Type: TypeInt}}
	results := Apply(env, rules)
	ok, failed := Summary(results)
	if ok != 1 || failed != 1 {
		t.Errorf("expected 1 ok 1 failed, got ok=%d failed=%d", ok, failed)
	}
}
