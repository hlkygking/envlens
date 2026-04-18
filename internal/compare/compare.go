package compare

import (
	"fmt"
	"sort"

	"github.com/user/envlens/internal/diff"
	"github.com/user/envlens/internal/parser"
)

// Result holds the full comparison between two env files.
type Result struct {
	BaseFile   string
	TargetFile string
	Entries    []diff.Entry
}

// Files parses two env files and returns a comparison Result.
func Files(basePath, targetPath string) (*Result, error) {
	base, err := parser.ParseFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("parsing base file %q: %w", basePath, err)
	}
	target, err := parser.ParseFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("parsing target file %q: %w", targetPath, err)
	}
	entries := diff.Compare(base, target)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return &Result{
		BaseFile:   basePath,
		TargetFile: targetPath,
		Entries:    entries,
	}, nil
}

// Summary returns counts per status.
func (r *Result) Summary() map[string]int {
	counts := map[string]int{}
	for _, e := range r.Entries {
		counts[string(e.Status)]++
	}
	return counts
}
