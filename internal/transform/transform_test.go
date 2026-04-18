package transform

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

func TestApply_Uppercase(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	res, err := Apply(env, []Rule{{Op: OpUppercase}})
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(res, "APP_ENV")
	if r == nil || r.Value != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %v", r)
	}
}

func TestApply_TrimSpace(t *testing.T) {
	env := map[string]string{"KEY": "  hello  "}
	res, err := Apply(env, []Rule{{Op: OpTrimSpace}})
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(res, "KEY")
	if r == nil || r.Value != "hello" {
		t.Errorf("expected 'hello', got %v", r)
	}
}

func TestApply_PrefixSuffix(t *testing.T) {
	env := map[string]string{"TOKEN": "abc"}
	res, err := Apply(env, []Rule{{Op: OpPrefix, Arg: "pre_"}, {Op: OpSuffix, Arg: "_suf"}})
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(res, "TOKEN")
	if r == nil || r.Value != "pre_abc_suf" {
		t.Errorf("unexpected value: %v", r)
	}
}

func TestApply_Replace(t *testing.T) {
	env := map[string]string{"URL": "http://localhost"}
	res, err := Apply(env, []Rule{{Op: OpReplace, Arg: "localhost", Arg2: "example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	r := findResult(res, "URL")
	if r == nil || r.Value != "http://example.com" {
		t.Errorf("unexpected value: %v", r)
	}
}

func TestApply_UnknownOp(t *testing.T) {
	env := map[string]string{"K": "v"}
	_, err := Apply(env, []Rule{{Op: Op("invalid")}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestSummary_Counts(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "WORLD"}
	res, _ := Apply(env, []Rule{{Op: OpUppercase}})
	transformed, unchanged := Summary(res)
	if transformed+unchanged != 2 {
		t.Errorf("expected 2 total, got %d+%d", transformed, unchanged)
	}
}
