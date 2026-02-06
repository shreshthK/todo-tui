package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shreshthkandari/todo-tui/internal/core"
	"github.com/shreshthkandari/todo-tui/internal/store"
	"github.com/shreshthkandari/todo-tui/internal/ui"
)

// main is the CLI entry point; it parses args, loads tasks, and dispatches commands.
func main() {
	// Parse the command-line arguments (excluding the program name).
	args := os.Args[1:]

	// Decide where the tasks JSON file should live on this machine.
	storePath, err := defaultStorePath()
	if err != nil {
		fmt.Printf("error: could not resolve tasks path: %v\n", err)
		os.Exit(1)
	}

	// Create the JSON-backed store using the resolved path.
	taskStore := store.NewJSONStore(storePath)

	// Load all tasks from disk before we execute any command.
	tasks, err := taskStore.Load()
	if err != nil {
		fmt.Printf("error: could not load tasks: %v\n", err)
		os.Exit(1)
	}

	// Dispatch to the right command implementation.
	if len(args) == 0 {
		if err := ui.StartTUI(taskStore, tasks); err != nil {
			fmt.Printf("error: could not start tui: %v\n", err)
			os.Exit(1)
		}
		return
	}

	switch args[0] {
	case "tui":
		if err := ui.StartTUI(taskStore, tasks); err != nil {
			fmt.Printf("error: could not start tui: %v\n", err)
			os.Exit(1)
		}
	case "add":
		handleAdd(taskStore, tasks, args[1:])
	case "list":
		handleList(tasks, args[1:])
	case "done":
		handleDone(taskStore, tasks, args[1:])
	case "delete":
		handleDelete(taskStore, tasks, args[1:])
	case "help":
		printUsage()
	default:
		fmt.Printf("error: unknown command %q\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

// defaultStorePath builds the full path to ~/.todo/tasks.json.
func defaultStorePath() (string, error) {
	// Ask the OS for the current user's home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Join the home directory with the plan's hidden folder path.
	return filepath.Join(home, ".todo", "tasks.json"), nil
}

// handleAdd creates a new task and saves it to disk.
func handleAdd(taskStore *store.JSONStore, tasks []core.Task, args []string) {
	// Join all remaining args into a single title string.
	title := strings.TrimSpace(strings.Join(args, " "))
	if title == "" {
		fmt.Println("error: missing task title")
		printUsage()
		os.Exit(1)
	}

	// Add the task to the in-memory list.
	updated, err := core.AddTask(tasks, title)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Persist the updated list to disk.
	if err := taskStore.Save(updated); err != nil {
		fmt.Printf("error: could not save tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("task added")
}

// handleList prints tasks to stdout, optionally filtered.
func handleList(tasks []core.Task, args []string) {
	// Use the first arg as a filter if provided; otherwise list all.
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	// Apply the filter in core so business rules stay centralized.
	filtered, err := core.ListTasks(tasks, filter)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Render the tasks in a friendly, readable format.
	printTasks(filtered)
}

// handleDone marks a task as completed and saves the result.
func handleDone(taskStore *store.JSONStore, tasks []core.Task, args []string) {
	// Parse the ID argument so we know which task to update.
	id, err := parseID(args)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	// Mark the task as done in memory.
	updated, err := core.MarkDone(tasks, id)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Persist the updated list to disk.
	if err := taskStore.Save(updated); err != nil {
		fmt.Printf("error: could not save tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("task marked done")
}

// handleDelete removes a task by ID and saves the result.
func handleDelete(taskStore *store.JSONStore, tasks []core.Task, args []string) {
	// Parse the ID argument so we know which task to remove.
	id, err := parseID(args)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	// Delete the task from the in-memory list.
	updated, err := core.DeleteTask(tasks, id)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// Persist the updated list to disk.
	if err := taskStore.Save(updated); err != nil {
		fmt.Printf("error: could not save tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("task deleted")
}

// parseID converts the first argument into an integer task ID.
func parseID(args []string) (int, error) {
	// Guard against missing arguments to avoid a panic.
	if len(args) == 0 {
		return 0, fmt.Errorf("missing id")
	}

	// Convert the string into an int.
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid id %q", args[0])
	}

	return id, nil
}

// printTasks renders a list of tasks in a simple terminal-friendly format.
func printTasks(tasks []core.Task) {
	// If there are no tasks, let the user know explicitly.
	if len(tasks) == 0 {
		fmt.Println("no tasks")
		return
	}

	// Print each task with a done marker and its ID.
	for _, t := range tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Printf("%s %d: %s\n", status, t.ID, t.Title)
	}
}

// printUsage explains the available commands and examples.
func printUsage() {
	fmt.Println("todo-tui (TUI + CLI)")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  todo-tui                  # start TUI")
	fmt.Println("  todo-tui tui              # start TUI")
	fmt.Println("  todo-tui add \"task title\"")
	fmt.Println("  todo-tui list [all|done|todo]")
	fmt.Println("  todo-tui done <id>")
	fmt.Println("  todo-tui delete <id>")
	fmt.Println("  todo-tui help")
}
