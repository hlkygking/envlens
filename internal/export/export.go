package export

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Format represents an export output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Export writes env map to a file in the given format.
func Export(env map[string]string, format Format, path string) error {
	var content string
	var err error
	switch format {
	case FormatDotenv:
		content = toDotenv(env)
	case FormatJSON:
		content, err = toJSON(env)
		if err != nil {
			return err
		}
	case FormatShell:
		content = toShell(env)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// Render returns the exported content as a string.
func Render(env map[string]string, format Format) (string, error) {
	switch format {
	case FormatDotenv:
		return toDotenv(env), nil
	case FormatJSON:
		return toJSON(env)
	case FormatShell:
		return toShell(env), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toDotenv(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "%s=%q\n", k, env[k])
	}
	return sb.String()
}

func toJSON(env map[string]string) (string, error) {
	b, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func toShell(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}
