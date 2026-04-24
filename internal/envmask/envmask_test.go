package envmask

import (
	"regexp"
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

func TestIsSensitive_BuiltIn(t *testing.T) {
	cases := []struct {
		key     string
		want    bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"SECRET", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, c := range cases {
		got := IsSensitive(c.key, nil)
		if got != c.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", c.key, got, c.want)
		}
	}
}

func TestIsSensitive_ExtraPattern(t *testing.T) {
	extra := []*regexp.Regexp{regexp.MustCompile(`(?i)magic`)}
	if !IsSensitive("MAGIC_VALUE", extra) {
		t.Error("expected MAGIC_VALUE to be sensitive with extra pattern")
	}
	if IsSensitive("MAGIC_VALUE", nil) {
		t.Error("expected MAGIC_VALUE NOT sensitive without extra pattern")
	}
}

func TestApply_RedactStrategy(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "myapp",
	}
	results := Apply(env, StrategyRedact, nil)

	r, ok := findResult(results, "DB_PASSWORD")
	if !ok {
		t.Fatal("DB_PASSWORD not in results")
	}
	if r.Masked != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", r.Masked)
	}
	if !r.WasMasked {
		t.Error("expected WasMasked=true")
	}

	r2, _ := findResult(results, "APP_NAME")
	if r2.Masked != "myapp" {
		t.Errorf("expected plain value, got %q", r2.Masked)
	}
}

func TestApply_PartialStrategy(t *testing.T) {
	env := map[string]string{"API_KEY": "abcdef"}
	results := Apply(env, StrategyPartial, nil)
	r, _ := findResult(results, "API_KEY")
	if r.Masked != "ab***" {
		t.Errorf("expected 'ab***', got %q", r.Masked)
	}
}

func TestApply_ShortValuePartial(t *testing.T) {
	env := map[string]string{"SECRET": "x"}
	results := Apply(env, StrategyPartial, nil)
	r, _ := findResult(results, "SECRET")
	if r.Masked != "***" {
		t.Errorf("expected '***' for short value, got %q", r.Masked)
	}
}

func TestToMap_ReturnsMasked(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret", "HOST": "localhost"}
	results := Apply(env, StrategyRedact, nil)
	m := ToMap(results)
	if m["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", m["DB_PASSWORD"])
	}
	if m["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", m["HOST"])
	}
}

func TestGetSummary_Counts(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_TOKEN":   "tok",
		"HOST":        "localhost",
		"PORT":        "5432",
	}
	results := Apply(env, StrategyRedact, nil)
	total, masked, plain := GetSummary(results)
	if total != 4 {
		t.Errorf("expected total=4, got %d", total)
	}
	if masked != 2 {
		t.Errorf("expected masked=2, got %d", masked)
	}
	if plain != 2 {
		t.Errorf("expected plain=2, got %d", plain)
	}
}
