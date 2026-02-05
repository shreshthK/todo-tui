package core

import (
	"errors"
	"time"
)

// Task represents a single todo item.
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

// AddTask appends a new task with a new ID and returns the updated slice.
func AddTask(tasks []Task, title string) ([]Task, error) {
	return nil, errors.New("AddTask: not implemented")
}

// DeleteTask removes a task by ID and returns the updated slice.
func DeleteTask(tasks []Task, id int) ([]Task, error) {
	return nil, errors.New("DeleteTask: not implemented")
}

// MarkDone marks a task as done by ID and returns the updated slice.
func MarkDone(tasks []Task, id int) ([]Task, error) {
	return nil, errors.New("MarkDone: not implemented")
}

// ListTasks returns tasks filtered by the provided filter string.
func ListTasks(tasks []Task, filter string) ([]Task, error) {
	return nil, errors.New("ListTasks: not implemented")
}
