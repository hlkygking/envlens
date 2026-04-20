// Package envscope provides functionality for scoping environment variables
// to named environments (e.g. dev, staging, prod) using prefix conventions
// or explicit mapping rules.
package envscope

import (
	"fmt"
	"strings"
)

// ScopeRule defines how to extract or assign a scope to an env key.
type ScopeRule struct {
	// Prefix is the key prefix that identifies this scope (e.g. "DEV_", "PROD_").
	Prefix string
	// Name is the logical scope name (e.g. "dev", "prod").
	Name string
	// StripPrefix controls whether the prefix is removed from the key in output.
	StripPrefix bool
}

// Result holds a single key scoped to one or more environments.
type Result struct {
	OriginalKey string
	ResolvedKey string
	Value       string
	Scope       string
	Matched     bool
}

// Summary contains aggregate statistics from a scope operation.
type Summary struct {
	Total    int
	Matched  int
	Unscoped int
	Scopes   map[string]int
}

// Apply scopes a flat map of environment variables using the provided rules.
// Keys matching a rule's prefix are tagged with the corresponding scope name.
// If StripPrefix is true, the prefix is removed from ResolvedKey.
// Keys not matching any rule are returned with Matched=false and Scope="unscoped".
func Apply(env map[string]string, rules []ScopeRule) []Result {
	results := make([]Result, 0, len(env))

	for key, value := range env {
		matched := false
		for _, rule := range rules {
			if rule.Prefix == "" {
				continue
			}
			if strings.HasPrefix(key, rule.Prefix) {
				resolvedKey := key
				if rule.StripPrefix {
					resolvedKey = strings.TrimPrefix(key, rule.Prefix)
				}
				results = append(results, Result{
					OriginalKey: key,
					ResolvedKey: resolvedKey,
					Value:       value,
					Scope:       rule.Name,
					Matched:     true,
				})
				matched = true
				break
			}
		}
		if !matched {
			results = append(results, Result{
				OriginalKey: key,
				ResolvedKey: key,
				Value:       value,
				Scope:       "unscoped",
				Matched:     false,
			})
		}
	}

	return results
}

// GetSummary computes aggregate counts from a slice of Results.
func GetSummary(results []Result) Summary {
	s := Summary{
		Total:  len(results),
		Scopes: make(map[string]int),
	}
	for _, r := range results {
		if r.Matched {
			s.Matched++
			s.Scopes[r.Scope]++
		} else {
			s.Unscoped++
		}
	}
	return s
}

// FilterByScope returns only results matching the given scope name.
func FilterByScope(results []Result, scope string) []Result {
	var out []Result
	for _, r := range results {
		if r.Scope == scope {
			out = append(out, r)
		}
	}
	return out
}

// ToMap converts a slice of Results into a key→value map using ResolvedKey.
// If multiple results share the same ResolvedKey (e.g. after prefix stripping),
// the last one wins and an error is returned listing the collisions.
func ToMap(results []Result) (map[string]string, error) {
	out := make(map[string]string, len(results))
	var collisions []string
	for _, r := range results {
		if _, exists := out[r.ResolvedKey]; exists {
			collisions = append(collisions, r.ResolvedKey)
		}
		out[r.ResolvedKey] = r.Value
	}
	if len(collisions) > 0 {
		return out, fmt.Errorf("key collisions after prefix stripping: %s", strings.Join(collisions, ", "))
	}
	return out, nil
}

// ParseRules parses a slice of strings in the format "PREFIX=scope[:strip]" into ScopeRules.
// Example: "DEV_=dev:strip" → ScopeRule{Prefix:"DEV_", Name:"dev", StripPrefix:true}
func ParseRules(raw []string) ([]ScopeRule, error) {
	rules := make([]ScopeRule, 0, len(raw))
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid scope rule %q: expected PREFIX=scope or PREFIX=scope:strip", s)
		}
		prefix := parts[0]
		rest := parts[1]
		strip := false
		if strings.HasSuffix(rest, ":strip") {
			strip = true
			rest = strings.TrimSuffix(rest, ":strip")
		}
		rules = append(rules, ScopeRule{
			Prefix:      prefix,
			Name:        rest,
			StripPrefix: strip,
		})
	}
	return rules, nil
}
