package diff_test

import (
	"testing"

	"github.com/envlens/envlens/internal/diff"
	"github.com/envlens/envlens/internal/parser"
	"github.com/stretchr/testify/assert"
)

func makeMap(pairs ...string) parser.EnvMap {
	m := make(parser.EnvMap)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func findEntry(entries []diff.Entry, key string) (diff.Entry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return diff.Entry{}, false
}

func TestCompare_Added(t *testing.T) {
	base := makeMap()
	head := makeMap("NEW_KEY", "value")
	r := diff.Compare(base, head)
	e, ok := findEntry(r.Entries, "NEW_KEY")
	assert.True(t, ok)
	assert.Equal(t, diff.Added, e.Kind)
	assert.Equal(t, "value", e.NewValue)
}

func TestCompare_Removed(t *testing.T) {
	base := makeMap("OLD_KEY", "val")
	head := makeMap()
	r := diff.Compare(base, head)
	e, ok := findEntry(r.Entries, "OLD_KEY")
	assert.True(t, ok)
	assert.Equal(t, diff.Removed, e.Kind)
	assert.Equal(t, "val", e.OldValue)
}

func TestCompare_Modified(t *testing.T) {
	base := makeMap("HOST", "localhost")
	head := makeMap("HOST", "prod.example.com")
	r := diff.Compare(base, head)
	e, ok := findEntry(r.Entries, "HOST")
	assert.True(t, ok)
	assert.Equal(t, diff.Modified, e.Kind)
	assert.Equal(t, "localhost", e.OldValue)
	assert.Equal(t, "prod.example.com", e.NewValue)
}

func TestCompare_Unchanged(t *testing.T) {
	base := makeMap("KEY", "same")
	head := makeMap("KEY", "same")
	r := diff.Compare(base, head)
	assert.False(t, r.HasChanges())
}

func TestHasChanges_True(t *testing.T) {
	base := makeMap("A", "1")
	head := makeMap("A", "2")
	r := diff.Compare(base, head)
	assert.True(t, r.HasChanges())
}
