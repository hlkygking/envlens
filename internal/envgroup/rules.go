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
