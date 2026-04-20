package envcount

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

func TestApply_Empty(t *testing.T) {
	env := map[string]string{"EMPTY_VAR": ""}
	results := Apply(env)
	r, ok := findResult(results, "EMPTY_VAR")
	if !ok {
		t.Fatal("expected result for EMPTY_VAR")
	}
	if r.Status != StatusEmpty {
		t.Errorf("expected empty, got %s", r.Status)
	}
	if r.Length != 0 {
		t.Errorf("expected length 0, got %d", r.Length)
	}
}

func TestApply_Short(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	results := Apply(env)
	r, ok := findResult(results, "PORT")
	if !ok {
		t.Fatal("expected result for PORT")
	}
	if r.Status != StatusShort {
		t.Errorf("expected short, got %s", r.Status)
	}
}

func TestApply_Medium(t *testing.T) {
	env := map[string]string{"APP_NAME": "my-awesome-application"}
	results := Apply(env)
	r, ok := findResult(results, "APP_NAME")
	if !ok {
		t.Fatal("expected result for APP_NAME")
	}
	if r.Status != StatusMedium {
		t.Errorf("expected medium, got %s", r.Status)
	}
}

func TestApply_Long(t *testing.T) {
	long := make([]byte, 100)
	for i := range long {
		long[i] = 'x'
	}
	env := map[string]string{"BIG_VAR": string(long)}
	results := Apply(env)
	r, ok := findResult(results, "BIG_VAR")
	if !ok {
		t.Fatal("expected result for BIG_VAR")
	}
	if r.Status != StatusLong {
		t.Errorf("expected long, got %s", r.Status)
	}
}

func TestApply_SensitiveKey(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "s3cr3t", "HOST": "localhost"}
	results := Apply(env)
	pass, _ := findResult(results, "DB_PASSWORD")
	host, _ := findResult(results, "HOST")
	if !pass.Sensitive {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
	if host.Sensitive {
		t.Error("expected HOST to not be sensitive")
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"EMPTY":    "",
		"SHORT":    "hi",
		"MEDIUM":   "hello-world-app",
		"API_KEY":  "short",
	}
	results := Apply(env)
	s := GetSummary(results)
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected 1 empty, got %d", s.Empty)
	}
	if s.Sensitive != 1 {
		t.Errorf("expected 1 sensitive, got %d", s.Sensitive)
	}
}
