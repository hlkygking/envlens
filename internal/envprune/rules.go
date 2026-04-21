package envprune

import (
	"fmt"
	"strings"
)

// Rule represents a single pruning directive parsed from a string.
type Rule struct {
	Type  string // "empty", "prefix", "pattern", "keep"
	Value string
}

// ParseRules parses a slice of rule strings into Rule values.
// Format: "empty", "prefix:DEBUG_", "pattern:^TMP_", "keep:APP_"
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if s == "empty" {
			rules = append(rules, Rule{Type: "empty"})
			continue
		}
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected type:value", s)
		}
		kind := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch kind {
		case "prefix", "pattern", "keep":
			rules = append(rules, Rule{Type: kind, Value: val})
		default:
			return nil, fmt.Errorf("unknown rule type %q", kind)
		}
	}
	return rules, nil
}

// ToOptions converts a slice of Rules into an Options struct.
func ToOptions(rules []Rule) Options {
	var opts Options
	for _, r := range rules {
		switch r.Type {
		case "empty":
			opts.RemoveEmpty = true
		case "prefix":
			opts.RemovePrefixes = append(opts.RemovePrefixes, r.Value)
		case "pattern":
			opts.RemovePatterns = append(opts.RemovePatterns, r.Value)
		case "keep":
			opts.KeepOnly = append(opts.KeepOnly, r.Value)
		}
	}
	return opts
}
