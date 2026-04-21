package envsplit

import (
	"regexp"
	"strings"
)

// Result holds the outcome of splitting a single env entry into buckets.
type Result struct {
	Key    string
	Value  string
	Bucket string // name of the matched bucket, or "default"
}

// Rule defines a named bucket and the key pattern that routes entries into it.
type Rule struct {
	Bucket  string
	Pattern string
}

// Summary holds aggregate counts per bucket.
type Summary struct {
	Buckets map[string]int
	Total   int
}

// Apply routes each key/value pair into a bucket based on the first matching rule.
// Entries that match no rule are placed in the "default" bucket.
func Apply(env map[string]string, rules []Rule) ([]Result, error) {
	compiled := make([]*regexp.Regexp, len(rules))
	for i, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		compiled[i] = re
	}

	results := make([]Result, 0, len(env))
	for k, v := range env {
		bucket := "default"
		for i, re := range compiled {
			if re.MatchString(k) {
				bucket = rules[i].Bucket
				break
			}
		}
		results = append(results, Result{Key: k, Value: v, Bucket: bucket})
	}
	return results, nil
}

// GetSummary returns per-bucket counts from a slice of results.
func GetSummary(results []Result) Summary {
	s := Summary{Buckets: make(map[string]int), Total: len(results)}
	for _, r := range results {
		s.Buckets[r.Bucket]++
	}
	return s
}

// FilterByBucket returns only results belonging to the named bucket.
func FilterByBucket(results []Result, bucket string) []Result {
	var out []Result
	for _, r := range results {
		if strings.EqualFold(r.Bucket, bucket) {
			out = append(out, r)
		}
	}
	return out
}

// ToMap converts a slice of results into a map of key→value.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}
