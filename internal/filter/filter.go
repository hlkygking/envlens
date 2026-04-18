package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	Prefix  string
	Suffix  string
	Pattern string
	Keys    []string
}

// Result holds a filtered key-value pair.
type Result struct {
	Key   string
	Value string
}

// Apply filters a map of env vars based on the provided Options.
// All non-empty criteria are ANDed together.
func Apply(env map[string]string, opts Options) ([]Result, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	var results []Result
	for k, v := range env {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, opts.Suffix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		if len(keySet) > 0 && !keySet[k] {
			continue
		}
		results = append(results, Result{Key: k, Value: v})
	}
	return results, nil
}

// Summary returns counts of total and matched entries.
type Summary struct {
	Total   int
	Matched int
}

func GetSummary(env map[string]string, results []Result) Summary {
	return Summary{Total: len(env), Matched: len(results)}
}
