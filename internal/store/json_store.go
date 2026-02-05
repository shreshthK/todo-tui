package store

import (
	"errors"

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
	return nil, errors.New("JSONStore.Load: not implemented")
}

// Save writes tasks to the JSON file.
func (s *JSONStore) Save(tasks []core.Task) error {
	return errors.New("JSONStore.Save: not implemented")
}
