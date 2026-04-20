package clone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your/envlens/internal/env"
)

// Result holds the outcome of a clone operation.
type Result struct {
	SourceFile string
	DestFile   string
	Keys       int
	Redacted   bool
}

// Options controls clone behaviour.
type Options struct {
	// RedactSensitive replaces sensitive values with empty string in the output.
	RedactSensitive bool
	// Extra holds additional key=value pairs to inject into the clone.
	Extra map[string]string
}

// Apply reads src, optionally redacts sensitive keys, merges Extra, and
// writes the result to dst. It returns a Result describing what happened.
func Apply(src, dst string, opts Options) (Result, error) {
	m, err := env.FromFile(src)
	if err != nil {
		return Result{}, fmt.Errorf("clone: read source: %w", err)
	}

	if opts.RedactSensitive {
		for k := range m {
			if isSensitive(k) {
				m[k] = ""
			}
		}
	}

	for k, v := range opts.Extra {
		m[k] = v
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return Result{}, fmt.Errorf("clone: mkdir: %w", err)
	}

	f, err := os.Create(dst)
	if err != nil {
		return Result{}, fmt.Errorf("clone: create dest: %w", err)
	}
	defer f.Close()

	for k, v := range m {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return Result{}, fmt.Errorf("clone: write: %w", err)
		}
	}

	return Result{
		SourceFile: src,
		DestFile:   dst,
		Keys:       len(m),
		Redacted:   opts.RedactSensitive,
	}, nil
}

// isSensitive reports whether key contains any well-known sensitive substrings.
// The comparison is case-insensitive so that keys like "db_password" and
// "DB_PASSWORD" are both matched.
func isSensitive(key string) bool {
	sensitive := []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE"}
	upper := strings.ToUpper(key)
	for _, s := range sensitive {
		if strings.Contains(upper, s) {
			return true
		}
	}
	return false
}
