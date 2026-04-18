package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single historical snapshot record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label"`
	Env       map[string]string `json:"env"`
}

// HistoryFile holds multiple snapshot entries.
type HistoryFile struct {
	Entries []Entry `json:"entries"`
}

// Append adds a new entry to the history file at the given path.
// If the file does not exist, it is created.
func Append(path, label string, env map[string]string) error {
	hf, err := loadFile(path)
	if err != nil {
		return fmt.Errorf("history load: %w", err)
	}
	hf.Entries = append(hf.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Env:       env,
	})
	return saveFile(path, hf)
}

// List returns all entries stored in the history file.
func List(path string) ([]Entry, error) {
	hf, err := loadFile(path)
	if err != nil {
		return nil, fmt.Errorf("history list: %w", err)
	}
	return hf.Entries, nil
}

// Latest returns the most recent entry, or an error if history is empty.
func Latest(path string) (Entry, error) {
	entries, err := List(path)
	if err != nil {
		return Entry{}, err
	}
	if len(entries) == 0 {
		return Entry{}, fmt.Errorf("history is empty")
	}
	return entries[len(entries)-1], nil
}

func loadFile(path string) (HistoryFile, error) {
	var hf HistoryFile
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return hf, nil
	}
	if err != nil {
		return hf, err
	}
	if err := json.Unmarshal(data, &hf); err != nil {
		return hf, fmt.Errorf("invalid history JSON: %w", err)
	}
	return hf, nil
}

func saveFile(path string, hf HistoryFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(hf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
