package envhash

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Entry holds the hash result for a single environment variable.
type Entry struct {
	Key   string
	Value string
	Hash  string
}

// Summary holds aggregate statistics from a hash operation.
type Summary struct {
	Total  int
	Hashed int
}

// Options controls how hashing is performed.
type Options struct {
	// KeysOnly hashes only the key names, ignoring values.
	KeysOnly bool
	// Salt is an optional prefix mixed into each hash.
	Salt string
}

// Apply hashes each key-value pair in env using SHA-256.
// Returns a slice of Entry with the resulting digests.
func Apply(env map[string]string, opts Options) []Entry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]Entry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		var input string
		if opts.KeysOnly {
			input = opts.Salt + k
		} else {
			input = opts.Salt + k + "=" + v
		}
		sum := sha256.Sum256([]byte(input))
		results = append(results, Entry{
			Key:   k,
			Value: v,
			Hash:  fmt.Sprintf("%x", sum),
		})
	}
	return results
}

// ToMap returns a map of key -> hash from the entry slice.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Hash
	}
	return m
}

// GetSummary returns aggregate counts for the given entries.
func GetSummary(entries []Entry) Summary {
	return Summary{
		Total:  len(entries),
		Hashed: len(entries),
	}
}

// CanonicalHash returns a single SHA-256 digest representing the entire
// environment map, suitable for change detection.
func CanonicalHash(env map[string]string, salt string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	if salt != "" {
		sb.WriteString(salt)
		sb.WriteByte('\n')
	}
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(env[k])
		sb.WriteByte('\n')
	}
	sum := sha256.Sum256([]byte(sb.String()))
	return fmt.Sprintf("%x", sum)
}
