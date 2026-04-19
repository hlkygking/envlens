package envtag

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

func TestApply_TagsMatchingKeys(t *testing.T) {
	rules, _ := ParseRules([]string{"^DB_:database", "^AWS_:cloud,infra"})
	env := map[string]string{"DB_HOST": "localhost", "AWS_REGION": "us-east-1", "PORT": "8080"}
	results := Apply(env, rules)

	r := findResult(results, "DB_HOST")
	if r == nil || !r.Tagged {
		t.Fatal("expected DB_HOST to be tagged")
	}
	if len(r.Tags) != 1 || r.Tags[0] != "database" {
		t.Errorf("unexpected tags: %v", r.Tags)
	}
}

func TestApply_MultipleTagsFromOneRule(t *testing.T) {
	rules, _ := ParseRules([]string{"^AWS_:cloud,infra"})
	env := map[string]string{"AWS_REGION": "us-east-1"}
	results := Apply(env, rules)

	r := findResult(results, "AWS_REGION")
	if r == nil {
		t.Fatal("missing AWS_REGION")
	}
	if len(r.Tags) != 2 {
		t.Errorf("expected 2 tags, got %v", r.Tags)
	}
}

func TestApply_UntaggedKey(t *testing.T) {
	rules, _ := ParseRules([]string{"^DB_:database"})
	env := map[string]string{"PORT": "8080"}
	results := Apply(env, rules)

	r := findResult(results, "PORT")
	if r == nil || r.Tagged {
		t.Fatal("expected PORT to be untagged")
	}
}

func TestParseRules_InvalidFormat(t *testing.T) {
	_, err := ParseRules([]string{"NOCOLON"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRules_InvalidRegex(t *testing.T) {
	_, err := ParseRules([]string{"[invalid:tag"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Tagged: true},
		{Key: "B", Tagged: false},
		{Key: "C", Tagged: true},
	}
	tagged, untagged := Summary(results)
	if tagged != 2 || untagged != 1 {
		t.Errorf("got tagged=%d untagged=%d", tagged, untagged)
	}
}
