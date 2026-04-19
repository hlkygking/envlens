package envset

import (
	"fmt"
	"sort"
	"strings"
)

// Entry represents a single environment variable.
type Entry struct {
	Key   string
	Value string
}

// Set holds a named collection of environment variables.
type Set struct {
	Name    string
	entries map[string]string
}

// New creates a new named Set from a map.
func New(name string, m map[string]string) *Set {
	copied := make(map[string]string, len(m))
	for k, v := range m {
		copied[k] = v
	}
	return &Set{Name: name, entries: copied}
}

// Get returns the value for a key and whether it exists.
func (s *Set) Get(key string) (string, bool) {
	v, ok := s.entries[key]
	return v, ok
}

// Set sets or updates a key.
func (s *Set) Set(key, value string) {
	s.entries[key] = value
}

// Delete removes a key from the set.
func (s *Set) Delete(key string) {
	delete(s.entries, key)
}

// Keys returns sorted keys.
func (s *Set) Keys() []string {
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ToMap returns a copy of the underlying map.
func (s *Set) ToMap() map[string]string {
	out := make(map[string]string, len(s.entries))
	for k, v := range s.entries {
		out[k] = v
	}
	return out
}

// Len returns the number of entries.
func (s *Set) Len() int {
	return len(s.entries)
}

// Contains reports whether the key exists.
func (s *Set) Contains(key string) bool {
	_, ok := s.entries[key]
	return ok
}

// String returns a human-readable summary.
func (s *Set) String() string {
	return fmt.Sprintf("Set(%s, %d keys)", s.Name, len(s.entries))
}

// Intersect returns keys present in both sets.
func Intersect(a, b *Set) []string {
	var common []string
	for _, k := range a.Keys() {
		if b.Contains(k) {
			common = append(common, k)
		}
	}
	return common
}

// Union returns all keys from both sets (sorted, deduplicated).
func Union(a, b *Set) []string {
	seen := make(map[string]struct{})
	for _, k := range a.Keys() {
		seen[k] = struct{}{}
	}
	for _, k := range b.Keys() {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Subtract returns keys in a but not in b.
func Subtract(a, b *Set) []string {
	var result []string
	for _, k := range a.Keys() {
		if !b.Contains(k) {
			result = append(result, k)
		}
	}
	return result
}

// HasPrefix returns all keys that start with prefix.
func (s *Set) HasPrefix(prefix string) []string {
	var result []string
	for _, k := range s.Keys() {
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}
	return result
}
