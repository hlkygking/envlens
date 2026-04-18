// Package parser reads .env files and returns key-value maps.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap is a map of environment variable key to value.
type EnvMap map[string]string

// ParseFile reads the file at path and returns an EnvMap.
// Lines starting with '#' or empty lines are ignored.
// Values may optionally be quoted with single or double quotes.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parser: %q line %d: %w", path, lineNum, err)
		}
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scan %q: %w", path, err)
	}

	return env, nil
}

func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line %q: missing '='")
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("invalid line %q: empty key", line)
	}

	value := strings.TrimSpace(parts[1])
	value = stripQuotes(value)

	return key, value, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
