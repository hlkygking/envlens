package envsample

import (
	"math/rand"
	"sort"
)

// Strategy controls how keys are selected during sampling.
type Strategy string

const (
	StrategyRandom Strategy = "random"
	StrategyFirst  Strategy = "first"
	StrategyLast   Strategy = "last"
)

// Result holds the outcome for a single key after sampling.
type Result struct {
	Key      string
	Value    string
	Included bool
}

// Options configures the sampling behaviour.
type Options struct {
	N        int
	Strategy Strategy
	Seed     int64
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		N:        10,
		Strategy: StrategyRandom,
		Seed:     42,
	}
}

// Apply samples up to N keys from env according to opts.
func Apply(env map[string]string, opts Options) []Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch opts.Strategy {
	case StrategyFirst:
		// already sorted ascending — take first N
	case StrategyLast:
		// reverse then take first N
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	default:
		// random shuffle with deterministic seed
		r := rand.New(rand.NewSource(opts.Seed))
		r.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	}

	selected := make(map[string]bool)
	limit := opts.N
	if limit > len(keys) {
		limit = len(keys)
	}
	for _, k := range keys[:limit] {
		selected[k] = true
	}

	// Rebuild results in stable alphabetical order.
	allKeys := make([]string, 0, len(env))
	for k := range env {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	results := make([]Result, 0, len(allKeys))
	for _, k := range allKeys {
		results = append(results, Result{
			Key:      k,
			Value:    env[k],
			Included: selected[k],
		})
	}
	return results
}

// GetSummary returns counts of included vs excluded keys.
func GetSummary(results []Result) (included, excluded int) {
	for _, r := range results {
		if r.Included {
			included++
		} else {
			excluded++
		}
	}
	return
}

// ToMap returns only the included key/value pairs.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string)
	for _, r := range results {
		if r.Included {
			m[r.Key] = r.Value
		}
	}
	return m
}
