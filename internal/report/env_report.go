package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// EnvSummary holds a summary of an env source for reporting.
type EnvSummary struct {
	Name  string            `json:"name"`
	Count int               `json:"count"`
	Keys  []string          `json:"keys"`
	Vars  map[string]string `json:"vars"`
}

// EnvTextReport writes a human-readable summary of an env source.
func EnvTextReport(w io.Writer, name string, vars map[string]string) {
	keys := sortedKeys(vars)
	fmt.Fprintf(w, "Source: %s (%d keys)\n", name, len(keys))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, vars[k])
	}
}

// EnvJSONReport writes a JSON summary of an env source.
func EnvJSONReport(w io.Writer, name string, vars map[string]string) error {
	keys := sortedKeys(vars)
	summary := EnvSummary{
		Name:  name,
		Count: len(vars),
		Keys:  keys,
		Vars:  vars,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(summary)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
