package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envlens/internal/mask"
)

// MaskedEntry represents a key/value pair with masking applied.
type MaskedEntry struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Masked  bool   `json:"masked"`
}

// MaskTextReport writes a human-readable masked env listing to w.
func MaskTextReport(env map[string]string, w io.Writer) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "=== Masked Environment ===")
	for _, k := range keys {
		v := env[k]
		if mask.IsSensitive(k) {
			v = mask.MaskValue(v)
			fmt.Fprintf(w, "  %s=%s [masked]\n", k, v)
		} else {
			fmt.Fprintf(w, "  %s=%s\n", k, v)
		}
	}
	fmt.Fprintf(w, "Total: %d keys\n", len(keys))
}

// MaskJSONReport writes a JSON array of MaskedEntry to w.
func MaskJSONReport(env map[string]string, w io.Writer) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]MaskedEntry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		isSens := mask.IsSensitive(k)
		if isSens {
			v = mask.MaskValue(v)
		}
		entries = append(entries, MaskedEntry{Key: k, Value: v, Masked: isSens})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
