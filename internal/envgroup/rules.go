package envgroup

import (
	"fmt"
	"strings"
)

// ParseRules parses rules from a slice of "name=pattern" strings.
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid rule %q: expected name=pattern", s)
		}
		rules = append(rules, Rule{Name: strings.TrimSpace(parts[0]), Pattern: strings.TrimSpace(parts[1])})
	}
	return rules, nil
}

// DefaultRules returns a set of common grouping rules.
func DefaultRules() []Rule {
	return []Rule{
		{Name: "database", Pattern: "^(DB|DATABASE|POSTGRES|MYSQL|REDIS)_"},
		{Name: "auth", Pattern: "^(AUTH|JWT|OAUTH|TOKEN|SECRET|API_KEY)"},
		{Name: "app", Pattern: "^APP_"},
		{Name: "logging", Pattern: "^LOG_"},
		{Name: "server", Pattern: "^(SERVER|HOST|PORT|ADDR)_?"},
	}
}

// MergeRules combines base rules with overrides. If a rule in overrides shares
// a name with a rule in base, the override takes precedence.
func MergeRules(base, overrides []Rule) []Rule {
	result := make([]Rule, 0, len(base)+len(overrides))
	seen := make(map[string]bool)
	for _, r := range overrides {
		seen[r.Name] = true
		result = append(result, r)
	}
	for _, r := range base {
		if !seen[r.Name] {
			result = append(result, r)
		}
	}
	return result
}
