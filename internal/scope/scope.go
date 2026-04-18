package scope

import "strings"

// Scope represents a named environment tier (e.g. dev, staging, prod).
type Scope struct {
	Name string
	Env  map[string]string
}

// Result holds the outcome of a scope comparison for a single key.
type Result struct {
	Key    string
	Values map[string]string // scope name -> value
	Uniform bool
}

// Compare takes multiple scopes and returns per-key consistency results.
func Compare(scopes []Scope) []Result {
	keySet := map[string]struct{}{}
	for _, s := range scopes {
		for k := range s.Env {
			keySet[k] = struct{}{}
		}
	}

	results := make([]Result, 0, len(keySet))
	for key := range keySet {
		values := map[string]string{}
		for _, s := range scopes {
			if v, ok := s.Env[key]; ok {
				values[s.Name] = v
			}
		}
		results = append(results, Result{
			Key:    key,
			Values: values,
			Uniform: isUniform(values),
		})
	}
	return results
}

// Summary returns counts of uniform vs divergent keys.
func Summary(results []Result) (uniform, divergent int) {
	for _, r := range results {
		if r.Uniform {
			uniform++
		} else {
			divergent++
		}
	}
	return
}

func isUniform(values map[string]string) bool {
	var first string
	var set bool
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if !strings.EqualFold(v, first) {
			return false
		}
	}
	return true
}
