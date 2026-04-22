package envlock

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// LockEntry represents a single locked key-value pair with its hash.
type LockEntry struct {
	Key       string `json:"key"`
	Hash      string `json:"hash"`
	LockedAt  string `json:"locked_at"`
}

// LockFile represents the full lock file structure.
type LockFile struct {
	Version   int         `json:"version"`
	CreatedAt string      `json:"created_at"`
	Entries   []LockEntry `json:"entries"`
}

// CheckResult represents the result of checking a key against a lock file.
type CheckResult struct {
	Key     string
	Status  string // "ok", "drifted", "missing"
	Expected string
	Actual   string
}

func hashValue(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

// Lock creates a LockFile from a map of env vars.
func Lock(env map[string]string) LockFile {
	now := time.Now().UTC().Format(time.RFC3339)
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]LockEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, LockEntry{
			Key:      k,
			Hash:     hashValue(env[k]),
			LockedAt: now,
		})
	}
	return LockFile{Version: 1, CreatedAt: now, Entries: entries}
}

// Save writes a LockFile to disk as JSON.
func Save(path string, lf LockFile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("envlock: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a LockFile from disk.
func Load(path string) (LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LockFile{}, fmt.Errorf("envlock: read: %w", err)
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return LockFile{}, fmt.Errorf("envlock: parse: %w", err)
	}
	return lf, nil
}

// Check compares a live env map against a LockFile and returns results.
func Check(lf LockFile, env map[string]string) []CheckResult {
	results := make([]CheckResult, 0, len(lf.Entries))
	for _, entry := range lf.Entries {
		val, ok := env[entry.Key]
		if !ok {
			results = append(results, CheckResult{
				Key:      entry.Key,
				Status:   "missing",
				Expected: entry.Hash,
			})
			continue
		}
		actual := hashValue(val)
		status := "ok"
		if actual != entry.Hash {
			status = "drifted"
		}
		results = append(results, CheckResult{
			Key:      entry.Key,
			Status:   status,
			Expected: entry.Hash,
			Actual:   actual,
		})
	}
	return results
}

// GetSummary returns counts of ok, drifted, and missing entries.
func GetSummary(results []CheckResult) (ok, drifted, missing int) {
	for _, r := range results {
		switch r.Status {
		case "ok":
			ok++
		case "drifted":
			drifted++
		case "missing":
			missing++
		}
	}
	return
}
