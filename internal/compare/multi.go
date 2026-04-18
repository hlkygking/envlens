package compare

import (
	"fmt"

	"github.com/user/envlens/internal/diff"
	"github.com/user/envlens/internal/parser"
)

// MultiResult holds pairwise comparisons for multiple files against a base.
type MultiResult struct {
	Base    string
	Targets []PairResult
}

// PairResult is one base-vs-target comparison.
type PairResult struct {
	Target  string
	Entries []diff.Entry
}

// MultiFiles compares a base env file against several target files.
func MultiFiles(basePath string, targetPaths []string) (*MultiResult, error) {
	base, err := parser.ParseFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base %q: %w", basePath, err)
	}
	result := &MultiResult{Base: basePath}
	for _, tp := range targetPaths {
		target, err := parser.ParseFile(tp)
		if err != nil {
			return nil, fmt.Errorf("parsing target %q: %w", tp, err)
		}
		entries := diff.Compare(base, target)
		result.Targets = append(result.Targets, PairResult{
			Target:  tp,
			Entries: entries,
		})
	}
	return result, nil
}

// AllKeys returns the union of all keys across all pair results.
func (m *MultiResult) AllKeys() []string {
	seen := map[string]struct{}{}
	for _, pr := range m.Targets {
		for _, e := range pr.Entries {
			seen[e.Key] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
