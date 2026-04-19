package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Type represents a target cast type.
type Type string

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
)

// Result holds the outcome of casting a single env value.
type Result struct {
	Key      string
	RawValue string
	CastType Type
	CastValue interface{}
	OK       bool
	Error    string
}

// Rule maps a key to a desired type.
type Rule struct {
	Key  string
	Type Type
}

// Apply casts env map values according to the provided rules.
func Apply(env map[string]string, rules []Rule) []Result {
	results := make([]Result, 0, len(rules))
	for _, rule := range rules {
		raw, exists := env[rule.Key]
		if !exists {
			results = append(results, Result{
				Key:      rule.Key,
				CastType: rule.Type,
				OK:       false,
				Error:    "key not found",
			})
			continue
		}
		res := cast(rule.Key, raw, rule.Type)
		results = append(results, res)
	}
	return results
}

func cast(key, raw string, t Type) Result {
	r := Result{Key: key, RawValue: raw, CastType: t}
	switch t {
	case TypeString:
		r.CastValue = raw
		r.OK = true
	case TypeInt:
		v, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			r.Error = fmt.Sprintf("cannot cast %q to int", raw)
		} else {
			r.CastValue = v
			r.OK = true
		}
	case TypeFloat:
		v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			r.Error = fmt.Sprintf("cannot cast %q to float", raw)
		} else {
			r.CastValue = v
			r.OK = true
		}
	case TypeBool:
		v, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			r.Error = fmt.Sprintf("cannot cast %q to bool", raw)
		} else {
			r.CastValue = v
			r.OK = true
		}
	default:
		r.Error = fmt.Sprintf("unknown type %q", t)
	}
	return r
}

// Summary returns counts of OK and failed casts.
func Summary(results []Result) (ok, failed int) {
	for _, r := range results {
		if r.OK {
			ok++
		} else {
			failed++
		}
	}
	return
}
