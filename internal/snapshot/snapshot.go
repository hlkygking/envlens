package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of an environment variable set.
type Snapshot struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path, label string, env map[string]string) error {
	s := Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Env:       env,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal error: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write error: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read error: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: parse error: %w", err)
	}
	return &s, nil
}

// Diff returns keys that were added, removed, or changed between two snapshots.
// It returns three maps: added, removed, and changed (keyed by variable name).
func Diff(base, other *Snapshot) (added, removed, changed map[string]string) {
	added = make(map[string]string)
	removed = make(map[string]string)
	changed = make(map[string]string)

	for k, v := range other.Env {
		if baseVal, ok := base.Env[k]; !ok {
			added[k] = v
		} else if baseVal != v {
			changed[k] = v
		}
	}
	for k, v := range base.Env {
		if _, ok := other.Env[k]; !ok {
			removed[k] = v
		}
	}
	return added, removed, changed
}

// Equal reports whether two snapshots contain identical environment variables.
func Equal(a, b *Snapshot) bool {
	if len(a.Env) != len(b.Env) {
		return false
	}
	for k, v := range a.Env {
		if bVal, ok := b.Env[k]; !ok || bVal != v {
			return false
		}
	}
	return true
}
