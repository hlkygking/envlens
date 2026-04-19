package envgroup

import (
	"fmt"
	"regexp"
	"sort"
)

// Group represents a named collection of env keys matching a pattern.
type Group struct {
	Name    string
	Pattern string
	Keys    []string
}

// Result holds all groups after applying rules to an env map.
type Result struct {
	Groups    []Group
	Ungrouped []string
}

// Rule defines a grouping rule.
type Rule struct {
	Name    string
	Pattern string
}

// Apply groups env keys according to the provided rules.
func Apply(env map[string]string, rules []Rule) (Result, error) {
	grouped := map[string]bool{}
	var groups []Group

	for _, rule := range rules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("invalid pattern %q: %w", rule.Pattern, err)
		}
		var matched []string
		for k := range env {
			if re.MatchString(k) {
				matched = append(matched, k)
				grouped[k] = true
			}
		}
		sort.Strings(matched)
		groups = append(groups, Group{Name: rule.Name, Pattern: rule.Pattern, Keys: matched})
	}

	var ungrouped []string
	for k := range env {
		if !grouped[k] {
			ungrouped = append(ungrouped, k)
		}
	}
	sort.Strings(ungrouped)

	return Result{Groups: groups, Ungrouped: ungrouped}, nil
}

// Summary returns counts per group and total ungrouped.
func Summary(r Result) map[string]int {
	sm := map[string]int{}
	for _, g := range r.Groups {
		sm[g.Name] = len(g.Keys)
	}
	sm["ungrouped"] = len(r.Ungrouped)
	return sm
}
