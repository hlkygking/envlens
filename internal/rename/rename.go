package rename

import "strings"

// Rule defines a rename operation for an env key.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of a single rename operation.
type Result struct {
	OldKey string
	NewKey string
	Value  string
	Renamed bool
	Skipped bool // key not found
}

// Apply renames keys in env according to rules.
// Keys not matching any rule are passed through unchanged.
func Apply(env map[string]string, rules []Rule) []Result {
	seen := make(map[string]bool)
	var results []Result

	for _, rule := range rules {
		val, ok := env[rule.From]
		if !ok {
			results = append(results, Result{OldKey: rule.From, NewKey: rule.To, Skipped: true})
			continue
		}
		seen[rule.From] = true
		results = append(results, Result{
			OldKey:  rule.From,
			NewKey:  rule.To,
			Value:   val,
			Renamed: true,
		})
	}

	for k, v := range env {
		if !seen[k] {
			results = append(results, Result{OldKey: k, NewKey: k, Value: v})
		}
	}
	return results
}

// ToMap converts results to a key/value map using new keys.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string)
	for _, r := range results {
		if !r.Skipped {
			out[r.NewKey] = r.Value
		}
	}
	return out
}

// Summary returns counts of renamed and skipped keys.
func Summary(results []Result) (renamed, skipped, unchanged int) {
	for _, r := range results {
		switch {
		case r.Skipped:
			skipped++
		case r.Renamed:
			renamed++
		default:
			unchanged++
		}
	}
	return
}

// ParseRules parses "OLD=NEW" strings into Rule slices.
func ParseRules(raw []string) []Rule {
	var rules []Rule
	for _, s := range raw {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 2 {
			rules = append(rules, Rule{From: strings.TrimSpace(parts[0]), To: strings.TrimSpace(parts[1])})
		}
	}
	return rules
}
