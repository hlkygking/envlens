package envalias

// Rule maps an old key name to a new alias key.
type Rule struct {
	From string
	To   string
}

// Result holds the outcome of applying an alias rule to a single key.
type Result struct {
	OriginalKey string
	AliasKey    string
	Value       string
	Applied     bool   // true if the source key existed
	Conflict    bool   // true if the alias key already existed in the env
}

// Summary holds aggregate counts from an Apply call.
type Summary struct {
	Total     int
	Applied   int
	Skipped   int
	Conflicts int
}

// Apply processes the given env map against the provided alias rules.
// For each rule, if the From key exists in env, a new entry is added under
// the To key (unless To already exists, which is flagged as a conflict).
// The original key is preserved in the returned map.
func Apply(env map[string]string, rules []Rule) ([]Result, map[string]string) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(rules))
	for _, r := range rules {
		val, exists := env[r.From]
		if !exists {
			results = append(results, Result{
				OriginalKey: r.From,
				AliasKey:    r.To,
				Applied:     false,
			})
			continue
		}
		_, conflict := out[r.To]
		if !conflict {
			out[r.To] = val
		}
		results = append(results, Result{
			OriginalKey: r.From,
			AliasKey:    r.To,
			Value:       val,
			Applied:     !conflict,
			Conflict:    conflict,
		})
	}
	return results, out
}

// GetSummary aggregates result counts.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch {
		case r.Conflict:
			s.Conflicts++
		case r.Applied:
			s.Applied++
		default:
			s.Skipped++
		}
	}
	return s
}

// ParseRules parses alias rules from a slice of "FROM=TO" strings.
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		for i := 0; i < len(s); i++ {
			if s[i] == '=' {
				from := s[:i]
				to := s[i+1:]
				if from != "" && to != "" {
					rules = append(rules, Rule{From: from, To: to})
				}
				break
			}
		}
	}
	return rules, nil
}
