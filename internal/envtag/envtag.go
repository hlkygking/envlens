package envtag

import (
	"fmt"
	"regexp"
	"strings"
)

// Tag represents a label attached to an env key.
type Tag struct {
	Key   string
	Value string
	Tags  []string
}

// Result holds tagging outcome for a single key.
type Result struct {
	Key    string
	Value  string
	Tags   []string
	Tagged bool
}

// Rule maps a pattern to a set of tags.
type Rule struct {
	Pattern *regexp.Regexp
	Tags    []string
}

// Apply tags env keys based on the provided rules.
func Apply(env map[string]string, rules []Rule) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		tags := []string{}
		for _, r := range rules {
			if r.Pattern.MatchString(k) {
				tags = append(tags, r.Tags...)
			}
		}
		results = append(results, Result{
			Key:    k,
			Value:  v,
			Tags:   unique(tags),
			Tagged: len(tags) > 0,
		})
	}
	return results
}

// ParseRules parses rules from strings like "PATTERN:tag1,tag2".
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected PATTERN:tags", s)
		}
		re, err := regexp.Compile(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", parts[0], err)
		}
		tags := strings.Split(parts[1], ",")
		rules = append(rules, Rule{Pattern: re, Tags: tags})
	}
	return rules, nil
}

// Summary returns counts of tagged vs untagged keys.
func Summary(results []Result) (tagged, untagged int) {
	for _, r := range results {
		if r.Tagged {
			tagged++
		} else {
			untagged++
		}
	}
	return
}

// FilterByTag returns only the results that contain the given tag.
func FilterByTag(results []Result, tag string) []Result {
	out := []Result{}
	for _, r := range results {
		for _, t := range r.Tags {
			if t == tag {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

func unique(tags []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, t := range tags {
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
