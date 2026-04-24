package transform

import (
	"fmt"
	"strings"
)

// ParseRules parses a slice of rule strings into Rule values.
// Format: "op" or "op:arg" or "op:arg:arg2"
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		r, err := parseRule(spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func parseRule(spec string) (Rule, error) {
	parts := strings.SplitN(spec, ":", 3)
	op := Op(strings.TrimSpace(parts[0]))
	switch op {
	case OpUppercase, OpLowercase, OpTrimSpace:
		return Rule{Op: op}, nil
	case OpPrefix, OpSuffix:
		if len(parts) < 2 {
			return Rule{}, fmt.Errorf("op %s requires an argument", op)
		}
		return Rule{Op: op, Arg: parts[1]}, nil
	case OpReplace:
		if len(parts) < 3 {
			return Rule{}, fmt.Errorf("op replace requires two arguments (format: replace:old:new)")
		}
		return Rule{Op: op, Arg: parts[1], Arg2: parts[2]}, nil
	default:
		if op == "" {
			return Rule{}, fmt.Errorf("op name must not be empty")
		}
		return Rule{}, fmt.Errorf("unknown op %q (supported: %s)", op, strings.Join(SupportedOps(), ", "))
	}
}

// SupportedOps returns a list of supported operation names.
func SupportedOps() []string {
	return []string{
		string(OpUppercase),
		string(OpLowercase),
		string(OpTrimSpace),
		string(OpPrefix),
		string(OpSuffix),
		string(OpReplace),
	}
}
