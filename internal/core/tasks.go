package core

import (
	"errors"
	"strings"
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
	// Normalize and validate the incoming title so we don't store empty tasks.
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		return nil, errors.New("AddTask: title cannot be empty")
	}

	// Find the maximum existing ID so the new task gets a stable, unique ID.
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	// Create the new task with a new ID and a timestamp for learning purposes.
	newTask := Task{
		ID:        maxID + 1,
		Title:     cleanTitle,
		Done:      false,
		CreatedAt: time.Now(),
	}

	// Append the task to the slice and return the updated list.
	tasks = append(tasks, newTask)
	return tasks, nil
}

// DeleteTask removes a task by ID and returns the updated slice.
func DeleteTask(tasks []Task, id int) ([]Task, error) {
	// Track whether we found the task so we can return a helpful error.
	found := false

	// Build a new slice that excludes the task with the matching ID.
	updated := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		updated = append(updated, t)
	}

	// If we never found the task, surface an error to the caller.
	if !found {
		return nil, errors.New("DeleteTask: task not found")
	}

	// Reassign IDs so they are contiguous after deletion.
	for i := range updated {
		updated[i].ID = i + 1
	}

	return updated, nil
}

// MarkDone marks a task as done by ID and returns the updated slice.
func MarkDone(tasks []Task, id int) ([]Task, error) {
	// Walk the slice and flip the Done flag on the matching task.
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = true
			return tasks, nil
		}
	}

	return nil, errors.New("MarkDone: task not found")
}

// ListTasks returns tasks filtered by the provided filter string.
func ListTasks(tasks []Task, filter string) ([]Task, error) {
	// Normalize the filter so we can accept variants like "Done" or "DONE".
	normalized := strings.TrimSpace(strings.ToLower(filter))

	// Treat empty or "all" as no filtering.
	if normalized == "" || normalized == "all" {
		return tasks, nil
	}

	// Filter the list based on the chosen mode.
	filtered := make([]Task, 0, len(tasks))
	switch normalized {
	case "done":
		for _, t := range tasks {
			if t.Done {
				filtered = append(filtered, t)
			}
		}
	case "todo", "pending":
		for _, t := range tasks {
			if !t.Done {
				filtered = append(filtered, t)
			}
		}
	default:
		return nil, errors.New("ListTasks: unknown filter (use all|done|todo)")
	}

	return filtered, nil
}
