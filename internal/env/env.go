package env

import (
	"fmt"
	"os"
	"strings"
)

// Source represents an environment variable source.
type Source struct {
	Name string
	Vars map[string]string
}

// FromFile loads a Source from a .env file using the parser.
func FromFile(path string) (*Source, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("env: read %s: %w", path, err)
	}
	vars := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = stripQuotes(val)
		vars[key] = val
	}
	return &Source{Name: path, Vars: vars}, nil
}

// FromMap creates a Source from an existing map.
func FromMap(name string, vars map[string]string) *Source {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Source{Name: name, Vars: copy}
}

// Merge combines multiple Sources, later sources override earlier ones.
func Merge(sources ...*Source) *Source {
	merged := make(map[string]string)
	for _, s := range sources {
		for k, v := range s.Vars {
			merged[k] = v
		}
	}
	name := "merged"
	if len(sources) > 0 {
		names := make([]string, len(sources))
		for i, s := range sources {
			names[i] = s.Name
		}
		name = strings.Join(names, "+")
	}
	return &Source{Name: name, Vars: merged}
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
