package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/shreshthkandari/todo-tui/internal/core"
)

// JSONStore loads and saves tasks from a JSON file on disk.
type JSONStore struct {
	Path string
}

// NewJSONStore constructs a JSONStore for the given path.
func NewJSONStore(path string) *JSONStore {
	return &JSONStore{Path: path}
}

// Load reads tasks from the JSON file.
func (s *JSONStore) Load() ([]core.Task, error) {
	// If the file does not exist yet, we treat that as "no tasks".
	if _, err := os.Stat(s.Path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []core.Task{}, nil
		}
		return nil, err
	}

	// Read the entire file into memory so we can decode it.
	raw, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}

	// An empty file should behave like an empty list.
	if len(raw) == 0 {
		return []core.Task{}, nil
	}

	// Decode JSON into the task slice.
	var tasks []core.Task
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Save writes tasks to the JSON file.
func (s *JSONStore) Save(tasks []core.Task) error {
	// Ensure the parent directory exists so the write doesn't fail.
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	// Encode tasks as pretty JSON so the file is human-readable.
	raw, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// Write the file atomically to avoid partially written files on crash.
	tmpPath := s.Path + ".tmp"
	if err := os.WriteFile(tmpPath, raw, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, s.Path); err != nil {
		return err
	}

	return nil
}
