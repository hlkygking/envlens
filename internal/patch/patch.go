package patch

import (
	"fmt"
	"strings"
)

// Op represents a patch operation type.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Rule defines a single patch instruction.
type Rule struct {
	Op    Op
	Key   string
	Value string // used by set
	To    string // used by rename
}

// Result holds the outcome of applying a patch rule.
type Result struct {
	Rule    Rule
	Applied bool
	Note    string
}

// Apply applies a list of patch rules to the given env map.
func Apply(env map[string]string, rules []Rule) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var results []Result
	for _, r := range rules {
		switch r.Op {
		case OpSet:
			out[r.Key] = r.Value
			results = append(results, Result{Rule: r, Applied: true, Note: "set"})
		case OpDelete:
			if _, ok := out[r.Key]; ok {
				delete(out, r.Key)
				results = append(results, Result{Rule: r, Applied: true, Note: "deleted"})
			} else {
				results = append(results, Result{Rule: r, Applied: false, Note: "key not found"})
			}
		case OpRename:
			if v, ok := out[r.Key]; ok {
				out[r.To] = v
				delete(out, r.Key)
				results = append(results, Result{Rule: r, Applied: true, Note: fmt.Sprintf("%s -> %s", r.Key, r.To)})
			} else {
				results = append(results, Result{Rule: r, Applied: false, Note: "key not found"})
			}
		default:
			results = append(results, Result{Rule: r, Applied: false, Note: "unknown op"})
		}
	}
	return out, results
}

// Summary returns counts of applied and skipped rules.
func Summary(results []Result) (applied, skipped int) {
	for _, r := range results {
		if r.Applied {
			applied++
		} else {
			skipped++
		}
	}
	return
}

// ParseRules parses rules from strings like "set:KEY=VALUE", "delete:KEY", "rename:OLD=NEW".
func ParseRules(strs []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range strs {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q", s)
		}
		op := Op(strings.TrimSpace(parts[0]))
		body := strings.TrimSpace(parts[1])
		switch op {
		case OpSet:
			kv := strings.SplitN(body, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("set rule requires KEY=VALUE: %q", s)
			}
			rules = append(rules, Rule{Op: op, Key: kv[0], Value: kv[1]})
		case OpDelete:
			rules = append(rules, Rule{Op: op, Key: body})
		case OpRename:
			kv := strings.SplitN(body, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("rename rule requires OLD=NEW: %q", s)
			}
			rules = append(rules, Rule{Op: op, Key: kv[0], To: kv[1]})
		default:
			return nil, fmt.Errorf("unknown op %q in rule %q", op, s)
		}
	}
	return rules, nil
}
