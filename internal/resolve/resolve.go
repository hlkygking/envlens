package resolve

import (
	"fmt"
	"os"
	"strings"
)

// Result holds the resolved value for a key.
type Result struct {
	Key      string
	Value    string
	Source   string // "file", "env", "default", "missing"
	Resolved bool
}

// Options controls resolution behaviour.
type Options struct {
	AllowEnvOverride bool
	Defaults         map[string]string
}

// Resolve takes a parsed env map and resolves each key against the
// process environment and provided defaults.
func Resolve(env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(env))

	seen := map[string]bool{}
	for k, v := range env {
		seen[k] = true
		results = append(results, resolveKey(k, v, opts))
	}

	// Include defaults that were not in the file.
	for k, v := range opts.Defaults {
		if seen[k] {
			continue
		}
		results = append(results, Result{
			Key:      k,
			Value:    v,
			Source:   "default",
			Resolved: true,
		})
	}

	return results
}

func resolveKey(key, fileVal string, opts Options) Result {
	if opts.AllowEnvOverride {
		if envVal, ok := os.LookupEnv(key); ok {
			return Result{Key: key, Value: envVal, Source: "env", Resolved: true}
		}
	}

	if fileVal != "" {
		return Result{Key: key, Value: fileVal, Source: "file", Resolved: true}
	}

	if def, ok := opts.Defaults[key]; ok {
		return Result{Key: key, Value: def, Source: "default", Resolved: true}
	}

	return Result{Key: key, Value: "", Source: "missing", Resolved: false}
}

// Summary returns a human-readable summary string.
func Summary(results []Result) string {
	total, resolved, missing := len(results), 0, 0
	for _, r := range results {
		if r.Resolved {
			resolved++
		} else {
			missing++
		}
	}
	return fmt.Sprintf("total=%d resolved=%d missing=%d",
		total, resolved, missing)
}

// MissingKeys returns keys that could not be resolved.
func MissingKeys(results []Result) []string {
	var out []string
	for _, r := range results {
		if !r.Resolved {
			out = append(out, r.Key)
		}
	}
	return out
}

// ToMap converts results to a plain key→value map (resolved only).
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		if r.Resolved {
			m[r.Key] = r.Value
		}
	}
	_ = strings.ToLower // keep import used
	return m
}
